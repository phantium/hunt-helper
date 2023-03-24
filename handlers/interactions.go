package handlers

import (
	"discordbot/internal/common"
	"discordbot/internal/orm"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	emoji "github.com/tmdvs/Go-Emoji-Utils"
	"golang.org/x/exp/slices"
)

var emojis = []string{
	"dragon",
	"kraken",
	"yeti",
	"maze",
	"abyssal",
}

func GetEmojis(s *discordgo.Session, names []string, guild_id string) map[string]string {
	emojis, _ := s.GuildEmojis(guild_id)
	emoji_result := map[string]string{}
	for _, emoji := range emojis {
		if slices.Contains(names, emoji.Name) {
			emoji_result[emoji.Name] = emoji.Name + ":" + emoji.ID
		}
	}
	return emoji_result
}

var game_request_minutes time.Duration = 15
var game_request_timeout time.Duration = 30

var (
	ComponentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"select_dungeon": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			fagtime := orm.GetFindAGame(i.Member.User.ID, i.GuildID)
			if fagtime.CreatedAt.IsZero() == false && !time.Now().After(fagtime.CreatedAt.Add(game_request_minutes*time.Minute)) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You can request a game every 10 minutes.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			} else {
				minValues := 1
				dungeon_components := []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.SelectMenu{
								CustomID:    "dungeon_finder",
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
											Name: "dragon",
											ID:   "1082313506700935199",
										},
									},
									{
										Label: "Kraken",
										Value: "Kraken",
										Emoji: discordgo.ComponentEmoji{
											Name: "kraken",
											ID:   "1082313504901578822",
										},
									},
									{
										Label: "Yeti",
										Value: "Yeti",
										Emoji: discordgo.ComponentEmoji{
											Name: "yeti",
											ID:   "1082333118729556038",
										},
									},
									{
										Label: "Maze",
										Value: "Maze",
										Emoji: discordgo.ComponentEmoji{
											Name: "maze",
											ID:   "1082313502208827422",
										},
									},
									{
										Label: "Abyssal",
										Value: "Abyssal",
										Emoji: discordgo.ComponentEmoji{
											Name: "abyssal",
											ID:   "1082313499922944000",
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
											Name: "dragon",
											ID:   "1082313506700935199",
										},
									},
									{
										Label: "Kraken",
										Value: "Kraken",
										Emoji: discordgo.ComponentEmoji{
											Name: "kraken",
											ID:   "1082313504901578822",
										},
									},
									{
										Label: "Yeti",
										Value: "Yeti",
										Emoji: discordgo.ComponentEmoji{
											Name: "yeti",
											ID:   "1082333118729556038",
										},
									},
									{
										Label: "Maze",
										Value: "Maze",
										Emoji: discordgo.ComponentEmoji{
											Name: "maze",
											ID:   "1082313502208827422",
										},
									},
									{
										Label: "Abyssal",
										Value: "Abyssal",
										Emoji: discordgo.ComponentEmoji{
											Name: "abyssal",
											ID:   "1082313499922944000",
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

		},
		"select_coop": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			guild := orm.GetGuildConfig(i.GuildID)
			fagtime := orm.GetFindAGame(i.Member.User.ID, i.GuildID)
			if fagtime.CreatedAt.IsZero() == false && !time.Now().After(fagtime.CreatedAt.Add(game_request_minutes*time.Minute)) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You can request a game every 10 minutes.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			} else {
				respond := &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Thank you, your co-op request has been posted!",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				}

				s.InteractionRespond(i.Interaction, respond)
				s.ChannelMessageSend(guild.ChannelBrowse, i.Member.Mention()+" wants a co-op run :id: "+orm.GetPlayerID(i.Member.User.ID))
				// orm.AddFindAGame(dg_msg.ID, i.ChannelID, i.GuildID, i.Member.User.ID, , "run")
			}

		},
		"select_event": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			guild := orm.GetGuildConfig(i.GuildID)
			fagtime := orm.GetFindAGame(i.Member.User.ID, i.GuildID)
			if fagtime.CreatedAt.IsZero() == false && !time.Now().After(fagtime.CreatedAt.Add(game_request_minutes*time.Minute)) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You can request a game every 10 minutes.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			} else {
				respond := &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Thank you, your event request has been posted!",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				}

				s.InteractionRespond(i.Interaction, respond)
				dg_msg, _ := s.ChannelMessageSend(guild.ChannelBrowse, i.Member.Mention()+" wants a team mate for this weeks Discord event :id: "+orm.GetPlayerID(i.Member.User.ID))
				orm.AddFindAGame(dg_msg.ID, i.ChannelID, i.GuildID, i.Member.User.ID, []string{"event"}, "run")
			}

		},

		"dungeon_finder": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			guild := orm.GetGuildConfig(i.GuildID)
			data := i.MessageComponentData()

			fagtime := orm.GetFindAGame(i.Member.User.ID, i.GuildID)

			if orm.GetPlayerID(i.Member.User.ID) == "" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry, but you need to set your Player ID first: <#" + guild.ChannelPlayerID + ">",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			} else if fagtime.CreatedAt.IsZero() == false && !time.Now().After(fagtime.CreatedAt.Add(game_request_minutes*time.Minute)) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You can request a find a game every 10 minutes.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			} else if guild.ChannelBrowse == "" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry, but the server admin needs to set the browse channel by using: !set channel browse",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			} else {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Thank you, your request has been posted!",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					panic(err)
				}

				g, _ := s.State.Guild(i.GuildID)
				final_role_ids := []string{}
				final_named_roles := []string{}
				for _, r := range data.Values {
					for _, rg := range g.Roles {
						if strings.ToLower(r) == strings.ToLower(emoji.RemoveAll(rg.Name)) {
							final_named_roles = append(final_named_roles, strings.ToLower(emoji.RemoveAll(rg.Name)))
							final_role_ids = append(final_role_ids, "<@&"+rg.ID+">")
						}
					}
				}
				dg_msg, _ := s.ChannelMessageSend(guild.ChannelBrowse, i.Member.Mention()+" wants a dungeon run for: "+strings.Join(final_role_ids, ", ")+" :id: "+orm.GetPlayerID(i.Member.User.ID))
				orm.AddFindAGame(dg_msg.ID, i.ChannelID, i.GuildID, i.Member.User.ID, final_named_roles, "run")
				emoji_list := GetEmojis(s, emojis, i.GuildID)
				for _, role := range final_named_roles {
					err := s.MessageReactionAdd(dg_msg.ChannelID, dg_msg.ID, emoji_list[role])
					if err != nil {
						log.Println(err)
					}
				}
				go func() {
					time.Sleep(game_request_timeout * time.Minute)
					orm.DeleteFindAGame(i.Member.User.ID, i.GuildID)
					s.ChannelMessageDelete(dg_msg.ChannelID, dg_msg.ID)
				}()
			}
		},
		"dungeon_finder_carry": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			guild := orm.GetGuildConfig(i.GuildID)
			data := i.MessageComponentData()

			fagtime := orm.GetFindAGame(i.Member.User.ID, i.GuildID)

			if orm.GetPlayerID(i.Member.User.ID) == "" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry, but you need to set your Player ID first: <#" + guild.ChannelPlayerID + ">",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			} else if fagtime.CreatedAt.IsZero() == false && !time.Now().After(fagtime.CreatedAt.Add(game_request_minutes*time.Minute)) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You can request a find a game every 10 minutes.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			} else if guild.ChannelBrowse == "" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry, but the server admin needs to set the browse channel by using: !set channel browse",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			} else {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Thank you, your request has been posted!",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					panic(err)
				}

				g, _ := s.State.Guild(i.GuildID)
				final_role_ids := []string{}
				final_named_roles := []string{}
				for _, r := range data.Values {
					for _, rg := range g.Roles {
						if strings.ToLower(r) == strings.ToLower(emoji.RemoveAll(rg.Name)) {
							final_named_roles = append(final_named_roles, strings.ToLower(emoji.RemoveAll(rg.Name)))
							final_role_ids = append(final_role_ids, "<@&"+rg.ID+">")
						}
					}
				}
				dg_msg, _ := s.ChannelMessageSend(guild.ChannelBrowse, i.Member.Mention()+" wants a dungeon carry for: "+strings.Join(final_role_ids, ", ")+" :id: "+orm.GetPlayerID(i.Member.User.ID))
				orm.AddFindAGame(dg_msg.ID, i.ChannelID, i.GuildID, i.Member.User.ID, final_named_roles, "carry")
				emoji_list := GetEmojis(s, emojis, i.GuildID)
				for _, role := range final_named_roles {
					err := s.MessageReactionAdd(dg_msg.ChannelID, dg_msg.ID, emoji_list[role])
					if err != nil {
						log.Println(err)
					}
				}
				go func() {
					time.Sleep(game_request_timeout * time.Minute)
					orm.DeleteFindAGame(i.Member.User.ID, i.GuildID)
					s.ChannelMessageDelete(dg_msg.ChannelID, dg_msg.ID)
				}()
			}
		},
	}
	CommandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"dungeon_finder": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if !common.MemberHasPermission(s, i.GuildID, i.Member.User.ID, discordgo.PermissionAdministrator) ||
				i.Member.User.ID == "76048145347252224" {

				respond := &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry, you are not allowed to use this command",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				}
				s.InteractionRespond(i.Interaction, respond)

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
		},
	}
)
