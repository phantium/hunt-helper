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
		&CommandPostPlayerID,
		&CommandRegisterPlayerID,
	}

	GlobalCommandHandlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"player_id":          CommandPlayerIDHandler,
		"player_id_reverse":  CommandPlayerIDReverseHandler,
		"player_id_post":     CommandPostPlayerIDHandler,
		"player_id_register": CommandRegisterPlayerIDHandler,
	}

	// guild
	commands = []*discordgo.ApplicationCommand{
		&CommandDungeonFinder,
		&CommandDungeonFinderRoles,
		&CommandHHChannels,
		&CommandDFTimeouts,
	}

	CommandHandlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"dungeon_finder": InteractionDungeonFinder,
		"df_roles":       CommandDungeonFinderRolesHandler,
		"hh_channels":    CommandHHChannelsHandler,
		"df_settings":    CommandDFTimeoutsHandler,
	}

	ComponentHandlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"dungeon_finder_run":           InteractionDungeonFinderRun,
		"dungeon_finder_carry":         InteractionDungeonFinderCarry,
		"dungeon_finder_provide_carry": InteractionDungeonFinderProvideCarry,
		// "select_dungeon":         InteractionSelectDungeon,
		"select_coop":            InteractionSelectCoop,
		"select_event":           InteractionSelectEvent,
		"dungeon_questions":      InteractionDungeonQuestions,
		"dungeons_carry":         InteractionDungeonsCarry,
		"dungeons_provide_carry": InteractionDungeonsProvideCarry,
		"dungeons_run":           InteractionDungeonsRun,
	}
)

func LoadGlobalCommands(s *discordgo.Session) {
	// log.Info("Loading global commands...")
	for _, command := range global_commands {
		// log.Infof("Loaded: %v", command.Name)
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", command)
		if err != nil {
			log.Warningf("Error loading global command: %s", err.Error())
		}
	}
}

func LoadGuildCommands(s *discordgo.Session, gc *discordgo.GuildCreate) {
	// log.Info("Loading guild commands...")
	for _, command := range commands {
		// log.Infof("Loaded: %v", command.Name)
		_, err := s.ApplicationCommandCreate(s.State.User.ID, gc.Guild.ID, command)
		if err != nil {
			log.Warningf("error loading guild command: %s", err.Error())
		}
	}
}

func UnloadGuildCommands(s *discordgo.Session, gc *discordgo.GuildCreate) {
	// cleanup previous slash commands
	commands, err := s.ApplicationCommands(s.State.User.ID, gc.Guild.ID)
	if err != nil {
		log.Warningf("unable to unregister guild command: %s", err.Error())
	}
	for _, c := range commands {
		if err := s.ApplicationCommandDelete(c.ApplicationID, gc.Guild.ID, c.ID); err != nil {
			log.Warningf("failed to remove %s command: %s", c.ID, err)
		}
	}
}
