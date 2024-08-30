package commands

import (
	"discordbot/internal/common"
	"discordbot/internal/orm"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var playerid_filter = regexp.MustCompile(`[A-Z]{8}`)

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

var CommandPlayerIDMenu = discordgo.ApplicationCommand{
	Name:        "player_id_menu",
	Description: "Hunt Royale Hunt Royale ID Menu",
}

// player_id_menu
func CommandMenuPlayerID(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !common.MemberHasPermission(s, i.GuildID, i.Member.User.ID, discordgo.PermissionAdministrator) {
		interactionResponseWithMessage(s, i, "Sorry, you are not allowed to use this command")

	} else {
		initial_components := []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Register your Player ID",
						CustomID: "player_id_registration",
						Style:    discordgo.PrimaryButton,
					},
				},
			},
		}

		respond := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content:    "Register your Hunt Royale Player ID",
				Components: initial_components,
			},
		}
		s.InteractionRespond(i.Interaction, respond)
	}
}

// player_id_registration
func InteractionCommandPlayerIDMenu(s *discordgo.Session, i *discordgo.InteractionCreate) {
	initial_components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					Label:     "Enter your Player ID!",
					CustomID:  "player_id",
					Style:     discordgo.TextInputShort,
					MinLength: 8,
					MaxLength: 8,
					Required:  true,
				},
			},
		},
	}

	respond := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Title:      "Player ID Registration",
			CustomID:   "player_id_registration",
			Components: initial_components,
		},
	}
	s.InteractionRespond(i.Interaction, respond)
	// }
}

func PlayerIDEntry(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()

	if data.CustomID != "player_id_registration" {
		return
	}

	data_input := strings.ToUpper(data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)

	// prevent repeat character input, not a valid player ID
	first_letter := string(data_input[0])
	if strings.Count(data_input, first_letter) == 8 {
		respond := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Your Player ID input is incorrect, try again.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		}
		s.InteractionRespond(i.Interaction, respond)
		return
	}

	player_id, _ := orm.GetMemberWithPlayerID(i.Member.User.ID)

	if !playerid_filter.MatchString(data_input) {
		respond := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Your Player ID input is incorrect, try again.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		}
		s.InteractionRespond(i.Interaction, respond)
		return
	}

	if playerid_filter.MatchString(data_input) && player_id.PlayerID == data_input {
		message := fmt.Sprintf("Thank you %s, but Hunt Royale :id: %s is already assigned to you", i.Member.User.String(), data_input)
		respond := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		}
		s.InteractionRespond(i.Interaction, respond)
	} else if player_id.PlayerID != "" {
		orm.DelMembersExistingPlayerID(i.Member.User.ID)
		orm.AddMemberWithPlayerID(i.Member.User.ID, data_input)
		message := fmt.Sprintf("Thank you %s, your Hunt Royale :id: has been updated to: %s", i.Member.User.String(), data_input)
		respond := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		}
		s.InteractionRespond(i.Interaction, respond)
	} else {
		message := fmt.Sprintf("Thank you %s, Hunt Royale :id: %s is now assigned to you", i.Member.User.String(), data_input)
		respond := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		}
		s.InteractionRespond(i.Interaction, respond)
		orm.AddMemberWithPlayerID(i.Member.User.ID, data_input)
	}
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
