package commands

import (
	"discordbot/internal/orm"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var minLength int = 8

var CommandPlayerIDReverse = discordgo.ApplicationCommand{
	Name:        "player_id_reverse",
	Description: "Hunt Royale get Discord member from Player ID",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "playerid",
			Required:    true,
			Description: "Hunt Royale Player ID",
			MinLength:   &minLength,
			MaxLength:   8,
		},
	},
}

func CommandPlayerIDReverseHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	var playerID string
	for _, opt := range options {
		if opt.Name == "playerid" {
			playerID = opt.StringValue()
			break
		}
	}

	memberID := orm.GetMemberID(playerID)

	var message string
	if memberID != "" {
		message = fmt.Sprintf("Hunt Royale :id: %s belongs to <@%s>", playerID, memberID)
	} else {
		message = fmt.Sprintf("Hunt Royale :id: %s is not registered to anyone", playerID)
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		return
	}
}
