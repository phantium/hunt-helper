package commands

import (
	"discordbot/internal/orm"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var CommandPostPlayerID = discordgo.ApplicationCommand{
	Name:        "player_id_post",
	Description: "Hunt Royale Post your Player ID",
	// Options: []*discordgo.ApplicationCommandOption{},
}

func CommandPostPlayerIDHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	playerID := orm.GetPlayerID(i.Member.User.ID)
	if playerID == "" {
		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You have no player ID registered! Use /player_id_register",
			},
		}

		err := s.InteractionRespond(i.Interaction, response)
		if err != nil {
			return
		}
		return
	}

	message_embed := []*discordgo.MessageEmbed{
		{
			Title:       fmt.Sprintf("%s's Hunt Royale :id:", i.Member.User.String()),
			Description: playerID,
			Color:       0x500497,
		},
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:          message_embed,
			AllowedMentions: &discordgo.MessageAllowedMentions{Users: []string{i.Member.User.ID}},
		},
	}

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		return
	}
}
