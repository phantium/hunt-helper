package commands

import (
	"discordbot/internal/orm"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var CommandRegisterPlayerID = discordgo.ApplicationCommand{
	Name:        "player_id_register",
	Description: "Hunt Royale register Hunt Royale ID",
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

func CommandRegisterPlayerIDHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	var newplayerID string
	for _, opt := range options {
		if opt.Name == "playerid" {
			newplayerID = opt.StringValue()
			break
		}
	}

	player_id, _ := orm.GetMemberWithPlayerID(i.Member.User.ID)

	var message string
	if player_id.PlayerID != "" && player_id.PlayerID == newplayerID {
		message = fmt.Sprintf("Thank you %s, but Hunt Royale :id: %s is already assigned to you", i.Member.User.String(), newplayerID)
	} else if player_id.PlayerID != "" {
		orm.DelMembersExistingPlayerID(i.Member.User.ID)
		orm.AddMemberWithPlayerID(i.Member.User.ID, newplayerID)
		message = fmt.Sprintf("Thank you %s, your Hunt Royale :id: has been updated to: %s", i.Member.User.String(), newplayerID)
	} else {
		message = fmt.Sprintf("Thank you %s, Hunt Royale :id: %s is now assigned to you", i.Member.User.String(), newplayerID)
		orm.AddMemberWithPlayerID(i.Member.User.ID, newplayerID)
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
