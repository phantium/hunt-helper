package commands

import (
	"discordbot/internal/common"
	"discordbot/internal/orm"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var CommandHHChannels = discordgo.ApplicationCommand{
	Name:        "hh_channels",
	Description: "Configure Hunt Helper Channels",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionString,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "Dungeon Finder Browse Channel",
					Value: "browse",
				},
				{
					Name:  "Dungeon Finder Board Channel",
					Value: "board",
				},
				{
					Name:  "PlayerID Channel",
					Value: "playerid",
				},
			},
			Name:        "channel-type",
			Required:    true,
			Description: "Channel type to configure",
		},
		{
			Type:        discordgo.ApplicationCommandOptionChannel,
			Name:        "channel-value",
			Required:    true,
			Description: "channel to use",
		},
	},
}

func CommandHHChannelsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	var channel_type string
	var channel_value *discordgo.Channel
	for _, opt := range options {
		switch opt.Name {
		case "channel-type":
			channel_type = opt.StringValue()
		case "channel-value":
			channel_value = opt.ChannelValue(s)
		}
	}

	orm.UpdateGuildConfig(i.GuildID, map[string]interface{}{"channel_" + channel_type: channel_value.ID})

	message := fmt.Sprintf("Configured channel type: %s with channel: %s", channel_type, channel_value.Mention())

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
