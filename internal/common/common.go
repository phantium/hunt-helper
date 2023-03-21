package common

import (
	"github.com/bwmarrin/discordgo"
)

func MemberHasPermission(s *discordgo.Session, guildID string, userID string, permission int64) bool {
	member, err := s.State.Member(guildID, userID)
	if err != nil {
		if member, err = s.GuildMember(guildID, userID); err != nil {
			return false
		}
	}

	// Iterate through the role IDs stored in member.Roles
	// to check permissions
	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			return false
		}
		if role.Permissions&permission != 0 {
			return true
		}
	}

	return false
}

func DeleteMessage(s *discordgo.Session, m *discordgo.Message) {
	s.ChannelMessageDelete(m.ChannelID, m.ID)
}
