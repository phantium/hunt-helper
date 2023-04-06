package commands

import (
	"discordbot/internal/common"
	"discordbot/internal/orm"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var CommandDFTimeouts = discordgo.ApplicationCommand{
	Name:        "df_settings",
	Description: "Configure Hunt Helper Dungeon Finder Settings",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionString,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "Dungeon Finder Request Delay Minutes (30-60)",
					Value: "fag_request_time",
				},
				{
					Name:  "Dungeon Finder Message Timeout Minutes (30-60)",
					Value: "fag_request_timeout",
				},
				{
					Name:  "Dungeon Finder Dungeon Select Limit (1-3)",
					Value: "fag_dungeon_select_limit",
				},
			},
			Name:        "setting-type",
			Required:    true,
			Description: "Setting Type",
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "setting-value",
			Required:    true,
			Description: "Setting Value",
		},
	},
}

func CommandDFTimeoutsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !common.MemberHasPermission(s, i.GuildID, i.Member.User.ID, discordgo.PermissionAdministrator) {
		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You do not have administrator permissions!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		}

		err := s.InteractionRespond(i.Interaction, response)
		if err != nil {
			return
		}
		return
	}
	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var timeout_type string
	var timeout_value int64

	if opt, ok := optionMap["setting-type"]; ok {
		timeout_type = opt.StringValue()
	}

	if opt, ok := optionMap["setting-value"]; ok {
		timeout_value = opt.IntValue()
	}

	switch timeout_type {
	case "fag_request_time":
		if timeout_value < 30 || timeout_value > 60 {
			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Value must be between 30 and 60",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			}

			err := s.InteractionRespond(i.Interaction, response)
			if err != nil {
				return
			}
			return
		}
	case "fag_request_timeout":
		if timeout_value < 30 || timeout_value > 60 {
			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Value must be between 30 and 60",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			}

			err := s.InteractionRespond(i.Interaction, response)
			if err != nil {
				return
			}
			return
		}
	case "fag_dungeon_select_limit":
		if timeout_value < 1 || timeout_value > 3 {
			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Value must be between 1 and 3",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			}

			err := s.InteractionRespond(i.Interaction, response)
			if err != nil {
				return
			}
			return
		}
	}

	orm.UpdateGuildConfig(i.GuildID, map[string]interface{}{timeout_type: fmt.Sprint(timeout_value)})

	message := fmt.Sprintf("%s configured with: %s", timeout_type, fmt.Sprint(timeout_value))

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
