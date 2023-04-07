package handlers

import (
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"

	"discordbot/internal/common"
	"discordbot/internal/orm"
)

var playerid_filter = regexp.MustCompile(`[A-Z]{8}`)

func ConfigurePlayerChannelID(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!set channel playerid" && common.MemberHasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionAdministrator) {
		orm.UpdateGuildConfig(m.GuildID, map[string]interface{}{"channel_player_id": m.ChannelID})
		s.ChannelMessageSend(m.ChannelID, "Channel Player ID set to: "+m.ChannelID+"!")
		log.Infof("Channel Player ID set to: %s on Guild ID: %s", m.ChannelID, m.GuildID)
	} else if m.Content == "!set channel playerid" && !common.MemberHasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionAdministrator) {
		s.ChannelMessageSend(m.ChannelID, "Sorry, you are not a server administrator!")
	}
}

func OnPlayerMessageID(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Do not show the notification on the #find-a-game channel
	if m.ChannelID != orm.GetGuildConfig(m.GuildID).ChannelPlayerID {
		return
	}

	var new_player_id string = strings.ToUpper(m.Content)

	if len(m.Content) == 8 && playerid_filter.MatchString(new_player_id) {
		player_id, _ := orm.GetMemberWithPlayerID(m.Author.ID)
		if player_id.PlayerID != "" && player_id.PlayerID == new_player_id {
			stored_msg, err := s.ChannelMessageSend(m.ChannelID, "Thank you "+m.Author.Mention()+", but I already have that player ID: **"+new_player_id+"** stored for you.")
			if err != nil {
				return
			}
			timer := time.NewTimer(30 * time.Second)
			go func() {
				<-timer.C
				err := s.ChannelMessageDelete(m.ChannelID, stored_msg.ID)
				if err != nil {
					return
				}
			}()
		} else if player_id.PlayerID != "" {
			orm.DelMembersExistingPlayerID(m.Author.ID)
			orm.AddMemberWithPlayerID(m.Author.ID, new_player_id)
			stored_msg, err := s.ChannelMessageSend(m.ChannelID, "Thank you "+m.Author.Mention()+", your player ID has been updated to: **"+new_player_id+"**, previous ID: *"+player_id.PlayerID+"*")
			if err != nil {
				return
			}
			timer := time.NewTimer(30 * time.Second)
			go func() {
				<-timer.C
				err := s.ChannelMessageDelete(m.ChannelID, stored_msg.ID)
				if err != nil {
					return
				}
			}()
		} else {
			orm.AddMemberWithPlayerID(m.Author.ID, new_player_id)
			stored_msg, err := s.ChannelMessageSend(m.ChannelID, "Thank you "+m.Author.Mention()+", your player ID: **"+new_player_id+"** has been stored.")
			if err != nil {
				return
			}
			timer := time.NewTimer(30 * time.Second)
			go func() {
				<-timer.C
				s.ChannelMessageDelete(m.ChannelID, stored_msg.ID)
			}()
		}
		// finally delete the original message from the user
		err := s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			return
		}
	}
}

func GetPlayerMessageIDBackLog(s *discordgo.Session) {
	log.Info("Starting player ID backlog processor")
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
	for _, guild := range guilds {
		guild_config := orm.GetGuildConfig(guild.ID)
		if guild_config.ChannelPlayerID == "" {
			// skip this guild
			continue
		}

		// Get the specified channel in the guild
		channel, err := s.State.Channel(guild_config.ChannelPlayerID)
		if err != nil {
			// handle error
			continue
		}

		// Skip the guild if the channel is not in it
		if channel.GuildID != guild.ID {
			continue
		}

		// Loop over each message in the channel
		var messages []*discordgo.Message
		var lastMessageID string
		for {
			partialMessages, err := s.ChannelMessages(channel.ID, 100, lastMessageID, "", "")
			if err != nil {
				return
			}
			messages = append(messages, partialMessages...)
			if len(partialMessages) < 100 {
				break
			}
			lastMessageID = partialMessages[len(partialMessages)-1].ID
		}

		for _, m := range messages {
			if m.Author.ID == s.State.User.ID {
				continue
			}

			var new_player_id string = strings.ToUpper(m.Content)

			if len(m.Content) == 8 && playerid_filter.MatchString(new_player_id) {
				player_id, _ := orm.GetMemberWithPlayerID(m.Author.ID)
				if player_id.PlayerID != "" && player_id.PlayerID == new_player_id {
					stored_msg, err := s.ChannelMessageSend(m.ChannelID, "Thank you "+m.Author.Mention()+", but I already have that player ID: **"+new_player_id+"** stored for you.")
					if err != nil {
						return
					}
					timer := time.NewTimer(30 * time.Second)
					go func() {
						<-timer.C
						err := s.ChannelMessageDelete(m.ChannelID, stored_msg.ID)
						if err != nil {
							return
						}
					}()
				} else if player_id.PlayerID != "" {
					orm.DelMembersExistingPlayerID(m.Author.ID)
					orm.AddMemberWithPlayerID(m.Author.ID, new_player_id)
					stored_msg, err := s.ChannelMessageSend(m.ChannelID, "Thank you "+m.Author.Mention()+", your player ID has been updated to: **"+new_player_id+"**, previous ID: *"+player_id.PlayerID+"*")
					if err != nil {
						return
					}
					timer := time.NewTimer(30 * time.Second)
					go func() {
						<-timer.C
						err := s.ChannelMessageDelete(m.ChannelID, stored_msg.ID)
						if err != nil {
							return
						}
					}()
				} else {
					orm.AddMemberWithPlayerID(m.Author.ID, new_player_id)
					stored_msg, err := s.ChannelMessageSend(m.ChannelID, "Thank you "+m.Author.Mention()+", your player ID: **"+new_player_id+"** has been stored.")
					if err != nil {
						return
					}
					timer := time.NewTimer(30 * time.Second)
					go func() {
						<-timer.C
						s.ChannelMessageDelete(m.ChannelID, stored_msg.ID)
					}()
				}
				// finally delete the original message from the user
				err := s.ChannelMessageDelete(m.ChannelID, m.ID)
				if err != nil {
					return
				}
			}
		}
	}
	log.Info("Done with player ID backlog processor")
}
