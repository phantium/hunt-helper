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

var channel_playerid string = "1082676228491853824"

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
			stored_msg, _ := s.ChannelMessageSend(m.ChannelID, "Thank you "+m.Author.Mention()+", but I already have that player ID: **"+new_player_id+"** stored for you.")
			go func() {
				time.Sleep(30 * time.Second)
				s.ChannelMessageDelete(m.ChannelID, stored_msg.ID)
			}()
		} else if player_id.PlayerID != "" {
			orm.DelMembersExistingPlayerID(m.Author.ID)
			orm.AddMemberWithPlayerID(m.Author.ID, new_player_id)
			stored_msg, _ := s.ChannelMessageSend(m.ChannelID, "Thank you "+m.Author.Mention()+", your player ID has been updated to: **"+new_player_id+"**, previous ID: *"+player_id.PlayerID+"*")
			go func() {
				time.Sleep(30 * time.Second)
				s.ChannelMessageDelete(m.ChannelID, stored_msg.ID)
			}()
		} else {
			orm.AddMemberWithPlayerID(m.Author.ID, new_player_id)
			stored_msg, _ := s.ChannelMessageSend(m.ChannelID, "Thank you "+m.Author.Mention()+", your player ID: **"+new_player_id+"** has been stored.")
			go func() {
				time.Sleep(30 * time.Second)
				s.ChannelMessageDelete(m.ChannelID, stored_msg.ID)
			}()
		}
		// finally delete the original message from the user
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}
