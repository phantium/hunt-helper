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
	"golang.org/x/exp/slices"
)

var cfg configuration.DiscordConfig

const discord_config string = "discord.yml"

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
	session.AddHandler(events.InteractionGlobalCreate)
	session.AddHandler(events.InteractionGuildCreate)
	session.AddHandler(events.Ready)

	defer session.Close()

	// Open a websocket connection to Discord and begin listening.
	must(session.Open())

	commands.LoadGlobalCommands(session)

	// capture reactions to messages
	session.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		guild := orm.GetGuildConfig(r.GuildID)
		allchannels := []string{}

		allchannels = append(allchannels, guild.ChannelDragon)
		allchannels = append(allchannels, guild.ChannelKraken)
		allchannels = append(allchannels, guild.ChannelYeti)
		allchannels = append(allchannels, guild.ChannelMaze)
		allchannels = append(allchannels, guild.ChannelAbyssal)
		allchannels = append(allchannels, guild.ChannelCoop)
		allchannels = append(allchannels, guild.ChannelEvent)

		if slices.Contains(allchannels, r.ChannelID) {
			handlers.FindAGameEmojiResponse(s, r)
		}
	})

	// cleanup messages if necessary
	commands.DungeonFinderIntegrityCheck(session)
	handlers.GetPlayerMessageIDBackLog(session)
	// handlers.FindAGameStatsPoster(session)

	// session.AddHandler(func(s *discordgo.Session, r *discordgo.MessageUpdate))

	// Wait for a termination signal from the operating system.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
