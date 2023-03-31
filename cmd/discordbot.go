package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"discordbot/commands"
	"discordbot/events"
	"discordbot/handlers"
	"discordbot/internal/configuration"
	"discordbot/internal/orm"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

var cfg configuration.DiscordConfig

const discord_config string = "discord.yml"

// var reaction_request_time time.Duration = 10

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

	// register find a game handlers
	session.AddHandler(handlers.OnPlayerMessageID)
	session.AddHandler(handlers.FindAGameStats)
	session.AddHandler(handlers.ConfigurePlayerChannelID)
	session.AddHandler(handlers.ConfigureBrowseChannelID)

	session.AddHandler(events.InteractionGlobalCreate)
	session.AddHandler(events.InteractionGuildCreate)
	session.AddHandler(events.Ready)

	defer session.Close()

	// Open a websocket connection to Discord and begin listening.
	must(session.Open())

	commands.LoadGlobalCommands(session)

	// session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 	switch i.Type {
	// 	case discordgo.InteractionApplicationCommand:
	// 		if h, ok := handlers.CommandsHandlers[i.ApplicationCommandData().Name]; ok {
	// 			h(s, i)
	// 		}
	// 	case discordgo.InteractionMessageComponent:
	// 		if h, ok := handlers.ComponentsHandlers[i.MessageComponentData().CustomID]; ok {
	// 			h(s, i)
	// 		}
	// 	}
	// })

	// capture reactions to messages
	session.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		guild := orm.GetGuildConfig(r.GuildID)
		if guild.ChannelBrowse == r.ChannelID {
			handlers.FindAGameEmojiResponse(s, r)
		}
	})

	// session.AddHandler(func(s *discordgo.Session, r *discordgo.MessageUpdate))

	// Wait for a termination signal from the operating system.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
