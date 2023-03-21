package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"

	"discordbot/handlers"
	"discordbot/internal/configuration"
	"discordbot/internal/orm"
	"discordbot/internal/slashcommands"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/ken/store"
)

var cfg configuration.DiscordConfig

const discord_config string = "discord.yml"

var reaction_request_time time.Duration = 10

func init() {
	rootCmd.AddCommand(DiscordBot)
	configuration.ReadConfig(&cfg, discord_config)
}

var DiscordBot = &cobra.Command{
	Use:   "bot",
	Short: "Discord Bot",
	Long:  "Discord Bot",
	Run: func(cmd *cobra.Command, args []string) {
		RunDiscordBot()
	},
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func RunDiscordBot() {
	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// register guild event handlers
	session.AddHandler(handlers.GuildCreate)
	session.AddHandler(handlers.GuildDelete)

	// // register find a game handlers
	session.AddHandler(handlers.OnPlayerMessageID)
	session.AddHandler(handlers.FindAGameStats)
	session.AddHandler(handlers.ConfigurePlayerChannelID)
	session.AddHandler(handlers.ConfigureBrowseChannelID)
	// session.AddHandler(handlers.ConfigureLFGChannelID)

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Open a websocket connection to Discord and begin listening.
	// err = dg.Open()
	// if err != nil {
	// 	fmt.Println("Error opening Discord connection: ", err)
	// 	return
	// }
	defer session.Close()

	k, err := ken.New(session, ken.Options{
		CommandStore: store.NewDefault(),
	})
	must(err)

	must(k.RegisterCommands(
		// new(slashcommands.BoardPost),
		new(slashcommands.CommandPlayerID),
	))

	defer k.Unregister()

	must(session.Open())

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := handlers.CommandsHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := handlers.ComponentsHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	// capture reactions to messages
	session.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		guild := orm.GetGuildConfig(r.GuildID)
		guild_info := orm.GetGuild(r.GuildID)
		// this means we are in the find a game channel, and we care what is being reacted
		if guild.ChannelBrowse == r.ChannelID {
			// skip if the bot is reacting
			if r.UserID == s.State.User.ID {
				return
			}

			// check if we care about this emoji at all
			allowedemojis := []string{"dragon", "kraken", "yeti", "maze", "abyssal"}
			if !slices.Contains(allowedemojis, r.Emoji.Name) {
				return
			}

			// check if user already responded to the message
			reaction := orm.GetFindAGameReaction(r.MessageID, r.UserID, r.GuildID)
			if reaction.CreatedAt.IsZero() == false && !time.Now().After(reaction.CreatedAt.Add(reaction_request_time*time.Minute)) {
				msgref := &discordgo.MessageReference{
					ChannelID: r.ChannelID,
					MessageID: r.MessageID,
					GuildID:   r.GuildID,
				}
				react_msg, _ := s.ChannelMessageSendReply(r.ChannelID, fmt.Sprintf("<@%s> you have already reacted to this game request!", r.UserID), msgref)
				go func() {
					time.Sleep(30 * time.Second)
					s.ChannelMessageDelete(r.ChannelID, react_msg.ID)
				}()
				return
			}
			if err != nil {
				log.Println(err)
			}
			// prevent user responding to their own request
			fag := orm.GetFindAGameByMsgID(r.MessageID, r.GuildID)
			if r.UserID == fag.UserID {
				msgref := &discordgo.MessageReference{
					ChannelID: r.ChannelID,
					MessageID: r.MessageID,
					GuildID:   r.GuildID,
				}
				react_msg, _ := s.ChannelMessageSendReply(r.ChannelID, fmt.Sprintf("<@%s> you cannot respond to your own request!", r.UserID), msgref)
				go func() {
					time.Sleep(30 * time.Second)
					s.ChannelMessageDelete(r.ChannelID, react_msg.ID)
				}()
				return
			}
			handlers.ReactToFindAGame(s, r.UserID, fag.UserID, guild_info, r.Emoji.Name)
			orm.AddFindAGameReaction(r.MessageID, r.GuildID, fag.UserID, r.Emoji.Name)
		}
	})

	// for _, guild := range *orm.GetGuilds() {
	// 	_, err = dg.ApplicationCommandCreate(dg.State.Application.ID, guild.GuildID, &discordgo.ApplicationCommand{
	// 		Name:        "dungeon_finder",
	// 		Description: "Post the interaction selecting dungeons",
	// 	})

	// 	if err != nil {
	// 		log.Fatalf("Cannot create slash command: %v", err)
	// 	}
	// }

	// _, err = dg.ApplicationCommandCreate(dg.State.Application.ID, guild_id, &discordgo.ApplicationCommand{
	// 	Name:        "dungeon_finder",
	// 	Description: "Post the interaction for dungeons",
	// })

	// if err != nil {
	// 	log.Fatalf("Cannot create slash command: %v", err)
	// }

	// Wait for a termination signal from the operating system.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
