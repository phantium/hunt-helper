package handlers

import (
	"discordbot/commands"
	"discordbot/internal/orm"

	"github.com/bwmarrin/discordgo"

	log "github.com/sirupsen/logrus"
)

// executed when added into a new server (guild)

func GuildCreate(s *discordgo.Session, gc *discordgo.GuildCreate) {
	guild := orm.GetGuild(gc.Guild.ID)
	if guild.GuildID == "" {
		orm.CreateGuild(&orm.Guilds{
			GuildID:   gc.Guild.ID,
			GuildName: gc.Guild.Name,
			OwnerID:   gc.Guild.OwnerID,
		})
		guildconfig := orm.GetGuildConfig(gc.Guild.ID)
		if guildconfig.GuildID == "" {
			orm.CreateGuildConfig(gc.Guild.ID)
		}
	}
	log.Infof("Joined server with guild ID: %s (Name: %s) (Owner: %s)!", gc.Guild.ID, gc.Guild.Name, gc.Guild.OwnerID)

	commands.LoadGuildCommands(s, gc)

}
