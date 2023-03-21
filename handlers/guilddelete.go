package handlers

import (
	"discordbot/internal/orm"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func GuildDelete(s *discordgo.Session, gd *discordgo.GuildDelete) {
	orm.DeleteGuild(gd.Guild.ID)
	log.Infof("Left server with guild ID: %s!", gd.Guild.ID)
}
