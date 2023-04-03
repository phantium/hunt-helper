package commands

import (
	"bytes"
	"discordbot/internal/common"
	"discordbot/internal/orm"
	"fmt"
	"log"
	"strings"
	"text/template"
	"time"

	"github.com/bwmarrin/discordgo"
	emoji "github.com/tmdvs/Go-Emoji-Utils"
)

var emojisv2 = map[string]string{
	"dragon":  "🐉",
	"kraken":  "🐙",
	"yeti":    "⛄",
	"maze":    "🏰",
	"abyssal": "😈",
	"event":   "💬",
}

var emojisv2_id = map[string]string{
	"dragon":  "1090337829499973632",
	"kraken":  "1090338068910837830",
	"yeti":    "1090338180328333462",
	"maze":    "1090338293708759140",
	"abyssal": "1090338381587820575",
	"event":   "1090338557140410429",
}

var CommandDungeonFinder = discordgo.ApplicationCommand{
	Name:        "dungeon_finder",
	Description: "Hunt Royale Dungeon Finder",
}

func deleteGameEntryAndMessage(s *discordgo.Session, i *discordgo.InteractionCreate, dg_msg *discordgo.Message) {
	orm.DeleteFindAGame(i.Member.User.ID, i.GuildID)
	err := s.ChannelMessageDelete(dg_msg.ChannelID, dg_msg.ID)
	if err != nil {
		return
	}
}

func deleteGameRequestAfterTimeout(s *discordgo.Session, i *discordgo.InteractionCreate, dg_msg *discordgo.Message, request_timeout time.Duration) {
	time.AfterFunc(request_timeout, func() {
		deleteGameEntryAndMessage(s, i, dg_msg)
	})
}

func interactionResponseWithMessage(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
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

func sendDungeonMessage(s *discordgo.Session, member *discordgo.Member, roles []*discordgo.Role, channelBrowse string, dungeonMessageTemplate string) (*discordgo.Message, error) {
	// Fill in the struct with the necessary variables
	messageData := struct {
		Member   string
		Roles    []*discordgo.Role
		PlayerID string
	}{
		Member:   member.Mention(),
		Roles:    roles,
		PlayerID: orm.GetPlayerID(member.User.ID),
	}

	// Create a new template and parse the message template
	tmpl, err := template.New("dungeonMessage").Parse(dungeonMessageTemplate)
	if err != nil {
		return nil, err
	}

	// Execute the template with the message data
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, messageData)
	if err != nil {
		return nil, err
	}

	// Send the message to the channel
	msg, err := s.ChannelMessageSend(channelBrowse, buf.String())
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// dungeon run "dungeon_finder_run"
func InteractionDungeonFinderRun(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild := orm.GetGuildConfig(i.GuildID)
	guild_config := orm.GetGuildConfig(i.GuildID)
	data := i.MessageComponentData()

	request_time := time.Duration(guild_config.FAGRequestTime) * time.Minute
	request_timeout := time.Duration(guild_config.FAGRequestTimeout) * time.Minute

	fagtime := orm.GetFindAGame(i.Member.User.ID, i.GuildID)

	if orm.GetPlayerID(i.Member.User.ID) == "" {
		interactionResponseWithMessage(s, i, "Sorry, but you need to set your Player ID first: <#"+guild.ChannelPlayerID+">")
	} else if !fagtime.CreatedAt.IsZero() && !time.Now().After(fagtime.CreatedAt.Add(request_time)) {
		interactionResponseWithMessage(s, i, fmt.Sprintf("You can request a game every: **%.2f minutes** wait: **%.2f minutes**", request_time.Minutes(), time.Since(fagtime.CreatedAt.Add(request_time)).Minutes()))
	} else if guild.ChannelBrowse == "" {
		interactionResponseWithMessage(s, i, "Sorry, but the server admin needs to set the browse channel by using: !set channel browse")
	} else {
		interactionResponseWithMessage(s, i, "Thank you, your request has been posted!")

		g, err := s.State.Guild(i.GuildID)
		if err != nil {
			return
		}
		final_role_ids := []*discordgo.Role{}
		final_named_roles := []string{}
		for _, r := range data.Values {
			for _, rg := range g.Roles {
				if strings.EqualFold(r, emoji.RemoveAll(rg.Name)) {
					final_named_roles = append(final_named_roles, strings.ToLower(emoji.RemoveAll(rg.Name)))
					// final_role_ids = append(final_role_ids, "<@&"+rg.ID+">")
					final_role_ids = append(final_role_ids, rg)
				}
			}
		}
		dg_msg, err := sendDungeonMessage(s, i.Member, final_role_ids, guild.ChannelBrowse, "**Run Request** - {{.Member}}: {{range $role := .Roles}}<@&{{$role.ID}}>{{end}} :id: {{.PlayerID}}")
		if err != nil {
			return
		}
		orm.AddFindAGame(dg_msg.ID, i.ChannelID, i.GuildID, i.Member.User.ID, final_named_roles, "run")
		for _, role := range final_named_roles {
			err := s.MessageReactionAdd(dg_msg.ChannelID, dg_msg.ID, emojisv2[role])
			if err != nil {
				log.Println(err)
			}
		}
		deleteGameRequestAfterTimeout(s, i, dg_msg, request_timeout)
	}
}

// dungeon carry "dungeon_finder_carry"
func InteractionDungeonFinderCarry(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild := orm.GetGuildConfig(i.GuildID)
	guild_config := orm.GetGuildConfig(i.GuildID)
	data := i.MessageComponentData()

	request_time := time.Duration(guild_config.FAGRequestTime) * time.Minute
	request_timeout := time.Duration(guild_config.FAGRequestTimeout) * time.Minute

	fagtime := orm.GetFindAGame(i.Member.User.ID, i.GuildID)

	if orm.GetPlayerID(i.Member.User.ID) == "" {
		interactionResponseWithMessage(s, i, "Sorry, but you need to set your Player ID first: <#"+guild.ChannelPlayerID+">")
	} else if !fagtime.CreatedAt.IsZero() && !time.Now().After(fagtime.CreatedAt.Add(request_time)) {
		interactionResponseWithMessage(s, i, fmt.Sprintf("You can request a game every: **%.2f minutes** wait: **%.2f minutes**", request_time.Minutes(), time.Since(fagtime.CreatedAt.Add(request_time)).Minutes()))
	} else if guild.ChannelBrowse == "" {
		interactionResponseWithMessage(s, i, "Sorry, but the server admin needs to set the browse channel by using: !set channel browse")
	} else {
		interactionResponseWithMessage(s, i, "Thank you, your request has been posted!")

		g, _ := s.State.Guild(i.GuildID)
		final_role_ids := []*discordgo.Role{}
		final_named_roles := []string{}
		for _, r := range data.Values {
			for _, rg := range g.Roles {
				if strings.EqualFold(r, emoji.RemoveAll(rg.Name)) {
					final_named_roles = append(final_named_roles, strings.ToLower(emoji.RemoveAll(rg.Name)))
					final_role_ids = append(final_role_ids, rg)
				}
			}
		}
		dg_msg, err := sendDungeonMessage(s, i.Member, final_role_ids, guild.ChannelBrowse, "**Carry Request** - {{.Member}}: {{range $role := .Roles}}<@&{{$role.ID}}>{{end}} :id: {{.PlayerID}}")
		if err != nil {
			return
		}
		orm.AddFindAGame(dg_msg.ID, i.ChannelID, i.GuildID, i.Member.User.ID, final_named_roles, "carry")
		// emoji_list := GetEmojis(s, emojis, i.GuildID)
		for _, role := range final_named_roles {
			err := s.MessageReactionAdd(dg_msg.ChannelID, dg_msg.ID, emojisv2[role])
			if err != nil {
				log.Println(err)
			}
		}
		deleteGameRequestAfterTimeout(s, i, dg_msg, request_timeout)
	}
}

// handle "dungeon_finder"
func InteractionDungeonFinder(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !common.MemberHasPermission(s, i.GuildID, i.Member.User.ID, discordgo.PermissionAdministrator) {

		interactionResponseWithMessage(s, i, "Sorry, you are not allowed to use this command")

	} else {
		initial_components := []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Dungeon",
						CustomID: "select_dungeon",
						Style:    discordgo.PrimaryButton,
					},
					discordgo.Button{
						Label:    "Co-op",
						CustomID: "select_coop",
						Style:    discordgo.DangerButton,
					},
					discordgo.Button{
						Label:    "Event",
						CustomID: "select_event",
						Style:    discordgo.SuccessButton,
					},
				},
			},
		}

		respond := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content:    "Dungeon, co-op and event game finder, choose your option.",
				Components: initial_components,
			},
		}
		s.InteractionRespond(i.Interaction, respond)
	}
}

// handle "select_dungeon"
func InteractionSelectDungeon(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild_config := orm.GetGuildConfig(i.GuildID)
	request_time := time.Duration(guild_config.FAGRequestTime) * time.Minute
	fagtime := orm.GetFindAGame(i.Member.User.ID, i.GuildID)

	if !fagtime.CreatedAt.IsZero() && !time.Now().After(fagtime.CreatedAt.Add(request_time)) {
		interactionResponseWithMessage(s, i, fmt.Sprintf("You can request a game every: **%.2f minutes** wait: **%.2f minutes**", request_time.Minutes(), time.Since(fagtime.CreatedAt.Add(request_time)).Minutes()))
	} else {
		minValues := 1
		dungeon_components := []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "dungeon_finder_run",
						Placeholder: "Choose your dungeons for a run:",
						// This is where confusion comes from. If you don't specify these things you will get single item select.
						// These fields control the minimum and maximum amount of selected items.
						MinValues: &minValues,
						MaxValues: 3,
						Options: []discordgo.SelectMenuOption{
							{
								Label: "Dragon",
								Value: "Dragon",
								// Default works the same for multi-select menus.
								// Default: false,
								Emoji: discordgo.ComponentEmoji{
									Name: emojisv2["dragon"],
								},
							},
							{
								Label: "Kraken",
								Value: "Kraken",
								Emoji: discordgo.ComponentEmoji{
									Name: emojisv2["kraken"],
								},
							},
							{
								Label: "Yeti",
								Value: "Yeti",
								Emoji: discordgo.ComponentEmoji{
									Name: emojisv2["yeti"],
								},
							},
							{
								Label: "Maze",
								Value: "Maze",
								Emoji: discordgo.ComponentEmoji{
									Name: emojisv2["maze"],
								},
							},
							{
								Label: "Abyssal",
								Value: "Abyssal",
								Emoji: discordgo.ComponentEmoji{
									Name: emojisv2["abyssal"],
								},
							},
						},
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "dungeon_finder_carry",
						Placeholder: "Choose your dungeons for a carry:",
						// This is where confusion comes from. If you don't specify these things you will get single item select.
						// These fields control the minimum and maximum amount of selected items.
						MinValues: &minValues,
						MaxValues: 3,
						Options: []discordgo.SelectMenuOption{
							{
								Label: "Dragon",
								Value: "Dragon",
								// Default works the same for multi-select menus.
								// Default: false,
								Emoji: discordgo.ComponentEmoji{
									Name: emojisv2["dragon"],
								},
							},
							{
								Label: "Kraken",
								Value: "Kraken",
								Emoji: discordgo.ComponentEmoji{
									Name: emojisv2["kraken"],
								},
							},
							{
								Label: "Yeti",
								Value: "Yeti",
								Emoji: discordgo.ComponentEmoji{
									Name: emojisv2["yeti"],
								},
							},
							{
								Label: "Maze",
								Value: "Maze",
								Emoji: discordgo.ComponentEmoji{
									Name: emojisv2["maze"],
								},
							},
							{
								Label: "Abyssal",
								Value: "Abyssal",
								Emoji: discordgo.ComponentEmoji{
									Name: emojisv2["abyssal"],
								},
							},
						},
					},
				},
			},
		}

		respond := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content:    "Select your dungeon(s) *up to 3* to run or be carried:",
				Components: dungeon_components,
				Flags:      discordgo.MessageFlagsEphemeral,
			},
		}
		s.InteractionRespond(i.Interaction, respond)
	}
}

// handle "select_coop"
func InteractionSelectCoop(s *discordgo.Session, i *discordgo.InteractionCreate) {
	fagtime := orm.GetFindAGame(i.Member.User.ID, i.GuildID)

	guild := orm.GetGuildConfig(i.GuildID)
	guild_config := orm.GetGuildConfig(i.GuildID)
	request_time := time.Duration(guild_config.FAGRequestTime) * time.Minute
	request_timeout := time.Duration(guild_config.FAGRequestTimeout) * time.Minute

	if !fagtime.CreatedAt.IsZero() && !time.Now().After(fagtime.CreatedAt.Add(request_time)) {
		interactionResponseWithMessage(s, i, fmt.Sprintf("You can request a game every: **%.2f minutes** wait: **%.2f minutes**", request_time.Minutes(), time.Since(fagtime.CreatedAt.Add(request_time)).Minutes()))
	} else {
		interactionResponseWithMessage(s, i, "Thank you, your co-op request has been posted!")
		dg_msg, err := sendDungeonMessage(s, i.Member, []*discordgo.Role{}, guild.ChannelBrowse, "**Co-Op Request** - {{.Member}} :id: {{.PlayerID}}")
		if err != nil {
			return
		}
		orm.AddFindAGame(dg_msg.ID, i.ChannelID, i.GuildID, i.Member.User.ID, []string{"coop"}, "run")
		deleteGameRequestAfterTimeout(s, i, dg_msg, request_timeout)
	}
}

// handle "select_event"
func InteractionSelectEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	fagtime := orm.GetFindAGame(i.Member.User.ID, i.GuildID)

	guild := orm.GetGuildConfig(i.GuildID)
	guild_config := orm.GetGuildConfig(i.GuildID)
	request_time := time.Duration(guild_config.FAGRequestTime) * time.Minute
	request_timeout := time.Duration(guild_config.FAGRequestTimeout) * time.Minute

	if !fagtime.CreatedAt.IsZero() && !time.Now().After(fagtime.CreatedAt.Add(request_time)) {
		interactionResponseWithMessage(s, i, fmt.Sprintf("You can request a game every: **%.2f minutes** wait: **%.2f minutes**", request_time.Minutes(), time.Since(fagtime.CreatedAt.Add(request_time)).Minutes()))
	} else {
		interactionResponseWithMessage(s, i, "Thank you, your event request has been posted!")
		dg_msg, err := sendDungeonMessage(s, i.Member, []*discordgo.Role{}, guild.ChannelBrowse, "**Weekly Event Request** - {{.Member}} :id: {{.PlayerID}}")
		if err != nil {
			return
		}
		orm.AddFindAGame(dg_msg.ID, i.ChannelID, i.GuildID, i.Member.User.ID, []string{"event"}, "run")
		deleteGameRequestAfterTimeout(s, i, dg_msg, request_timeout)
	}
}
