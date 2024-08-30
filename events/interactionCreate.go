package events

import (
	"github.com/bwmarrin/discordgo"

	"discordbot/commands"
)

func InteractionGlobalCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if cmd, ok := commands.GlobalCommandHandlers[i.ApplicationCommandData().Name]; ok {
			cmd(s, i)
		}
	}
}

func InteractionGuildCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if cmd, ok := commands.CommandHandlers[i.ApplicationCommandData().Name]; ok {
			cmd(s, i)
		}
	case discordgo.InteractionMessageComponent:
		if cmp, ok := commands.ComponentHandlers[i.MessageComponentData().CustomID]; ok {
			cmp(s, i)
		}
	case discordgo.InteractionModalSubmit:
		if mdl, ok := commands.ModalHandlers[i.ModalSubmitData().CustomID]; ok {
			mdl(s, i)
		}
	}
}
