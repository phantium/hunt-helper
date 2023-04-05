package commands

import (
	"bytes"
	"discordbot/internal/common"
	"discordbot/internal/orm"
	"fmt"
	"strings"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
	emoji "github.com/tmdvs/Go-Emoji-Utils"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/exp/slices"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var emojisv2 = map[string]string{
	"dragon":  "🐉",
	"kraken":  "🐙",
	"yeti":    "⛄",
	"maze":    "🏰",
	"abyssal": "😈",
	"event":   "💬",
	"coop":    "👥",
}

// var emojisv2_id = map[string]string{
// 	"dragon":  "1090337829499973632",
// 	"kraken":  "1090338068910837830",
// 	"yeti":    "1090338180328333462",
// 	"maze":    "1090338293708759140",
// 	"abyssal": "1090338381587820575",
// 	"event":   "1090338557140410429",
// }

var CommandDungeonFinder = discordgo.ApplicationCommand{
	Name:        "dungeon_finder",
	Description: "Hunt Royale Dungeon Finder",
}

var CommandDungeonFinderRoles = discordgo.ApplicationCommand{
	Name:        "df_roles",
	Description: "Configure the roles required for dungeon_finder to function.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionString,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "dragon",
					Value: "dragon",
				},
				{
					Name:  "kraken",
					Value: "kraken",
				},
				{
					Name:  "yeti",
					Value: "yeti",
				},
				{
					Name:  "maze",
					Value: "maze",
				},
				{
					Name:  "abyssal",
					Value: "abyssal",
				},
				{
					Name:  "coop",
					Value: "coop",
				},
				{
					Name:  "event",
					Value: "event",
				},
			},
			Name:        "name",
			Required:    true,
			Description: "one of dragon,kraken,yeti,maze,abyssal,event,coop",
		},
		{
			Type:        discordgo.ApplicationCommandOptionRole,
			Name:        "role",
			Required:    true,
			Description: "role to use",
		},
	},
}

func CommandDungeonFinderRolesHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	var name string
	var role *discordgo.Role
	for _, opt := range options {
		switch opt.Name {
		case "name":
			name = opt.StringValue()
		case "role":
			role = opt.RoleValue(s, i.GuildID)
		}
	}

	orm.UpdateGuildConfig(i.GuildID, map[string]interface{}{"role_" + name: role.ID})

	message := fmt.Sprintf("Configured role: %s with the server role: %s", name, role.Mention())

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

func sendDungeonMessage(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.MessageComponentInteractionData, gameType string, dungeonMessageTemplate string) {
	// Fetch guilds using pagination
	var guilds []*discordgo.Guild
	var lastGuildID string
	for {
		partialGuilds, err := s.UserGuilds(100, lastGuildID, "")
		if err != nil {
			// handle error
			return
		}
		for _, g := range partialGuilds {
			guilds = append(guilds, &discordgo.Guild{
				ID: g.ID,
			})
		}
		if len(partialGuilds) < 100 {
			break
		}
		lastGuildID = partialGuilds[len(partialGuilds)-1].ID
	}

	for _, guild := range guilds {
		guild_config := orm.GetGuildConfig(guild.ID)
		request_timeout := time.Duration(guild_config.FAGRequestTimeout) * time.Minute

		if guild_config.ChannelBrowse == "" {
			continue
		}

		// roles magic v2
		roles := []*discordgo.Role{}
		named_roles := []string{}
		roles_mapping := map[string]string{
			"dragon":  guild_config.RoleDragon,
			"kraken":  guild_config.RoleKraken,
			"yeti":    guild_config.RoleYeti,
			"maze":    guild_config.RoleMaze,
			"abyssal": guild_config.RoleAbyssal,
			"coop":    guild_config.RoleCoop,
			"event":   guild_config.RoleEvent,
		}

		// loop over roles
		for name, id := range roles_mapping {
			fetch_role, err := s.State.Role(guild.ID, id)
			if err != nil {
				continue
			}
			if slices.Contains(data.Values, cases.Title(language.Und, cases.NoLower).String(name)) {
				roles = append(roles, fetch_role)
				named_roles = append(named_roles, emoji.RemoveAll(strings.ToLower(fetch_role.Name)))
			}
		}

		// get the guild name of the original request
		origin_guild, err := s.State.Guild(i.GuildID)
		if err != nil {
			continue
		}
		// Fill in the struct with the necessary variables
		messageData := struct {
			Member   string
			Guild    string
			Roles    []*discordgo.Role
			PlayerID string
		}{
			Member:   i.Member.Mention(),
			Guild:    origin_guild.Name,
			Roles:    roles,
			PlayerID: orm.GetPlayerID(i.Member.User.ID),
		}

		// Create a new template and parse the message template
		tmpl, err := template.New("dungeonMessage").Parse(dungeonMessageTemplate)
		if err != nil {
			return
		}

		// Execute the template with the message data
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, messageData)
		if err != nil {
			return
		}

		role_excluded_gametypes := []string{"coop", "event"}
		if !slices.Contains(role_excluded_gametypes, gameType) {
			// if we don't have any roles for the server, just continue
			if len(roles) == 0 {
				continue
			}
			// if we don't match the requested available roles, just continue
			if len(data.Values) != len(roles) {
				continue
			}
		}

		// Send the message to the channel
		msg, err := s.ChannelMessageSend(guild_config.ChannelBrowse, buf.String())
		if err != nil {
			continue
		}
		orm.AddFindAGame(msg.ID, i.ChannelID, i.GuildID, i.Member.User.ID, named_roles, gameType)
		deleteGameRequestAfterTimeout(s, i, msg, request_timeout)

		// reaction emoji roles
		for _, gametype := range data.Values {
			err := s.MessageReactionAdd(msg.ChannelID, msg.ID, emojisv2[strings.ToLower(gametype)])
			if err != nil {
				log.Println(err)
			}
		}

		// reaction emoji special roles
		special_gametypes := []string{"coop", "event"}
		if slices.Contains(special_gametypes, gameType) {
			err := s.MessageReactionAdd(msg.ChannelID, msg.ID, emojisv2[gameType])
			if err != nil {
				log.Println(err)
			}
		}

	}
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

// dungeon run "dungeon_finder_run"
func InteractionDungeonFinderRun(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild := orm.GetGuildConfig(i.GuildID)
	guild_config := orm.GetGuildConfig(i.GuildID)
	data := i.MessageComponentData()

	request_time := time.Duration(guild_config.FAGRequestTime) * time.Minute
	// request_timeout := time.Duration(guild_config.FAGRequestTimeout) * time.Minute

	fagtime := orm.GetFindAGame(i.Member.User.ID)

	if orm.GetPlayerID(i.Member.User.ID) == "" {
		interactionResponseWithMessage(s, i, "Sorry, but you need to set your Player ID first: <#"+guild.ChannelPlayerID+">")
	} else if !fagtime.CreatedAt.IsZero() && !time.Now().After(fagtime.CreatedAt.Add(request_time)) {
		interactionResponseWithMessage(s, i, fmt.Sprintf("You can request a game every: **%.2f minutes** wait: **%.2f minutes**", request_time.Minutes(), time.Since(fagtime.CreatedAt.Add(request_time)).Minutes()))
	} else if guild.ChannelBrowse == "" {
		interactionResponseWithMessage(s, i, "Sorry, but the server admin needs to set the browse channel by using: !set channel browse")
	} else {
		interactionResponseWithMessage(s, i, "Thank you, your request has been posted!")

		sendDungeonMessage(s, i, data, "run", "**Run Request** - {{.Member}} @ {{.Guild}}: {{range $role := .Roles}}<@&{{$role.ID}}> {{end}} :id: {{.PlayerID}}")
	}
}

// dungeon carry "dungeon_finder_carry"
func InteractionDungeonFinderCarry(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild := orm.GetGuildConfig(i.GuildID)
	guild_config := orm.GetGuildConfig(i.GuildID)
	data := i.MessageComponentData()

	request_time := time.Duration(guild_config.FAGRequestTime) * time.Minute
	// request_timeout := time.Duration(guild_config.FAGRequestTimeout) * time.Minute

	fagtime := orm.GetFindAGame(i.Member.User.ID)

	if orm.GetPlayerID(i.Member.User.ID) == "" {
		interactionResponseWithMessage(s, i, "Sorry, but you need to set your Player ID first: <#"+guild.ChannelPlayerID+">")
	} else if !fagtime.CreatedAt.IsZero() && !time.Now().After(fagtime.CreatedAt.Add(request_time)) {
		interactionResponseWithMessage(s, i, fmt.Sprintf("You can request a game every: **%.2f minutes** wait: **%.2f minutes**", request_time.Minutes(), time.Since(fagtime.CreatedAt.Add(request_time)).Minutes()))
	} else if guild.ChannelBrowse == "" {
		interactionResponseWithMessage(s, i, "Sorry, but the server admin needs to set the browse channel by using: !set channel browse")
	} else {
		interactionResponseWithMessage(s, i, "Thank you, your request has been posted!")

		sendDungeonMessage(s, i, data, "carry", "**Carry Request** - {{.Member}} @ {{.Guild}}: {{range $role := .Roles}}<@&{{$role.ID}}>{{end}} :id: {{.PlayerID}}")

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
	fagtime := orm.GetFindAGame(i.Member.User.ID)

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
	fagtime := orm.GetFindAGame(i.Member.User.ID)

	// guild := orm.GetGuildConfig(i.GuildID)
	guild_config := orm.GetGuildConfig(i.GuildID)
	request_time := time.Duration(guild_config.FAGRequestTime) * time.Minute
	// request_timeout := time.Duration(guild_config.FAGRequestTimeout) * time.Minute

	if !fagtime.CreatedAt.IsZero() && !time.Now().After(fagtime.CreatedAt.Add(request_time)) {
		interactionResponseWithMessage(s, i, fmt.Sprintf("You can request a game every: **%.2f minutes** wait: **%.2f minutes**", request_time.Minutes(), time.Since(fagtime.CreatedAt.Add(request_time)).Minutes()))
	} else {
		interactionResponseWithMessage(s, i, "Thank you, your co-op request has been posted!")
		sendDungeonMessage(s, i, discordgo.MessageComponentInteractionData{Values: []string{"coop"}}, "coop", "**Co-Op Request** - {{.Member}} @ {{.Guild}} :id: {{.PlayerID}}")

	}
}

// handle "select_event"
func InteractionSelectEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	fagtime := orm.GetFindAGame(i.Member.User.ID)

	// guild := orm.GetGuildConfig(i.GuildID)
	guild_config := orm.GetGuildConfig(i.GuildID)
	request_time := time.Duration(guild_config.FAGRequestTime) * time.Minute
	// request_timeout := time.Duration(guild_config.FAGRequestTimeout) * time.Minute

	if !fagtime.CreatedAt.IsZero() && !time.Now().After(fagtime.CreatedAt.Add(request_time)) {
		interactionResponseWithMessage(s, i, fmt.Sprintf("You can request a game every: **%.2f minutes** wait: **%.2f minutes**", request_time.Minutes(), time.Since(fagtime.CreatedAt.Add(request_time)).Minutes()))
	} else {
		interactionResponseWithMessage(s, i, "Thank you, your event request has been posted!")
		sendDungeonMessage(s, i, discordgo.MessageComponentInteractionData{Values: []string{"event"}}, "event", "**Weekly Event Request** - {{.Member}} @ {{.Guild}} :id: {{.PlayerID}}")

	}
}

// during startup check integrity of the messages still left in the database
// remove the messages leftover from the guilds after validating the messages exist
func DungeonFinderIntegrityCheck(s *discordgo.Session) {
	// Get all guilds the bot is a part of
	var guilds []*discordgo.UserGuild
	var lastGuildID string
	for {
		partialGuilds, err := s.UserGuilds(100, lastGuildID, "")
		if err != nil {
			// handle error
			return
		}
		guilds = append(guilds, partialGuilds...)
		if len(partialGuilds) < 100 {
			break
		}
		lastGuildID = partialGuilds[len(partialGuilds)-1].ID
	}

	// Loop over each guild
	log.Info("Starting Dungeon Finder integrity check.. hang tight!")
	for _, guild := range guilds {
		guild_config := orm.GetGuildConfig(guild.ID)
		if guild_config.ChannelBrowse == "" {
			continue
		}

		// Get the specified channel in the guild
		channel, err := s.State.Channel(guild_config.ChannelBrowse)
		if err != nil {
			// handle error
			continue
		}

		// Skip the guild if the channel is not in it
		if channel.GuildID != guild.ID {
			continue
		}

		// Get the bot's messages in the channel
		var messages []*discordgo.Message
		var lastMessageID string
		for {
			partialMessages, err := s.ChannelMessages(channel.ID, 100, lastMessageID, "", "")
			if err != nil {
				// handle error
				continue
			}
			messages = append(messages, partialMessages...)
			if len(partialMessages) < 100 {
				break
			}
			lastMessageID = partialMessages[len(partialMessages)-1].ID
		}

		apology_text := "Hello! I was just restarted.\nI'm cleaning up LFG messages no longer tied to an interaction.\nPlease repost if your LFG was deleted.\nThank you!\n\n*`this message self destructs in 2 minutes`*"
		apology_msg, err := s.ChannelMessageSend(guild_config.ChannelBrowse, apology_text)
		if err != nil {
			// did the channel disappear?
			continue
		}
		time.AfterFunc(2*time.Minute, func() {
			s.ChannelMessageDelete(guild_config.ChannelBrowse, apology_msg.ID)
		})

		// Loop over each message
		for _, message := range messages {
			// Delete the message if it was created by the bot
			if message.Author.ID == s.State.User.ID {
				// delete the message
				err := s.ChannelMessageDelete(channel.ID, message.ID)
				if err != nil {
					continue
				}

				// check if it still exists
				_, err = orm.GetFindAGameByMsgID(message.ID)
				if err != nil {
					continue
				} else {
					// delete from orm
					orm.DeleteFindAGameByMessageID(message.ID)
				}
			}
		}
	}
	log.Info("Done with Dungeon Finder cleanup")
}
