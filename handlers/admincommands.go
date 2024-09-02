package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const AdminUserID = "76048145347252224"

// LeaveServer handles the !leave command to make the bot leave a specified server
func LeaveServer(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Content, "!leave ") {
		return
	}

	args := strings.Split(m.Content, " ")
	if len(args) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !leave <server_id>")
		return
	}

	// Check if the user is the admin
	if m.Author.ID != AdminUserID {
		s.ChannelMessageSend(m.ChannelID, "You don't have permission to use this command.")
		return
	}

	serverID := args[1]
	err := s.GuildLeave(serverID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Failed to leave server: %v", err))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully left server with ID: %s", serverID))
}

// ListServers handles the !servers command to list all servers the bot is in
func ListServers(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Content, "!servers") {
		return
	}

	// Check if the user is the admin
	if m.Author.ID != AdminUserID {
		s.ChannelMessageSend(m.ChannelID, "You don't have permission to use this command.")
		return
	}

	guilds := s.State.Guilds
	response := "Servers I'm currently in:\n"

	for _, guild := range guilds {
		response += fmt.Sprintf("- %s (ID: %s)\n", guild.Name, guild.ID)
	}

	s.ChannelMessageSend(m.ChannelID, response)
}
