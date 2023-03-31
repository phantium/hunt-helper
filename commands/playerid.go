package commands

import (
	"discordbot/internal/orm"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var CommandPlayerID = discordgo.ApplicationCommand{
	Name:        "player_id",
	Description: "Hunt Royale get users Hunt Royale ID",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Required:    true,
			Description: "@ a discord user",
		},
	},
}

func CommandPlayerIDHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var playerid_message string
	var user *discordgo.User

	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	if opt, ok := optionMap["user"]; ok {
		user = opt.UserValue(s)
	}

	player_id := orm.GetPlayerID(user.ID)
	if player_id != "" {
		playerid_message = fmt.Sprintf("<@%s> Hunt Royale :id: %s", user.ID, player_id)
	} else {
		playerid_message = fmt.Sprintf("<@%s> has no Hunt Royale :id: registered", user.ID)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: playerid_message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
