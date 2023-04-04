package commands

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	// global
	global_commands = []*discordgo.ApplicationCommand{
		&CommandPlayerID,
		&CommandPlayerIDReverse,
	}

	GlobalCommandHandlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"player_id":         CommandPlayerIDHandler,
		"player_id_reverse": CommandPlayerIDReverseHandler,
	}

	// guild
	commands = []*discordgo.ApplicationCommand{
		&CommandDungeonFinder,
		&CommandDungeonFinderRoles,
		&CommandHHChannels,
	}

	CommandHandlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"dungeon_finder": InteractionDungeonFinder,
		"df_roles":       CommandDungeonFinderRolesHandler,
		"hh_channels":    CommandHHChannelsHandler,
	}

	ComponentHandlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"dungeon_finder_run":   InteractionDungeonFinderRun,
		"dungeon_finder_carry": InteractionDungeonFinderCarry,
		"select_dungeon":       InteractionSelectDungeon,
		"select_coop":          InteractionSelectCoop,
		"select_event":         InteractionSelectEvent,
	}
)

func LoadGlobalCommands(s *discordgo.Session) {
	log.Info("Loading global commands...")
	for _, command := range global_commands {
		// log.Infof("Loaded: %v", command.Name)
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", command)
		if err != nil {
			log.Warningf("Error loading global command: %s", err.Error())
		}
	}
}

func LoadGuildCommands(s *discordgo.Session, gc *discordgo.GuildCreate) {
	log.Info("Loading guild commands...")
	for _, command := range commands {
		// log.Infof("Loaded: %v", command.Name)
		_, err := s.ApplicationCommandCreate(s.State.User.ID, gc.Guild.ID, command)
		if err != nil {
			log.Warningf("Error loading guild command: %s", err.Error())
		}
	}
}
