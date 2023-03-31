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
	var memberid_message string
	var playerid string

	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	if opt, ok := optionMap["playerid"]; ok {
		playerid = opt.StringValue()
	}

	member_id := orm.GetMemberID(playerid)
	if member_id != "" {
		memberid_message = fmt.Sprintf("Hunt Royale :id: %s belongs to <@%s>", playerid, member_id)
	} else {
		memberid_message = fmt.Sprintf("Hunt Royale :id: %s is not registered to anyone", playerid)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: memberid_message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
