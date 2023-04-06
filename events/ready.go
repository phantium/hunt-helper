package events

import (
	"discordbot/utils"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

const (
	statusIntervalPeriod = 60 * time.Minute
)

var (
	statusTexts = []string{
		"Hunt Royale",
		"hide and seek",
	}
)

func Ready(session *discordgo.Session, event *discordgo.Ready) {
	err := session.UpdateGameStatus(0, utils.RandomString(statusTexts))
	if err != nil {
		log.Warningf("Unable to set status: %s", err.Error())
	}

	ticker := time.NewTicker(statusIntervalPeriod)
	for range ticker.C {
		err := session.UpdateGameStatus(0, utils.RandomString(statusTexts))
		if err != nil {
			log.Warningf("An error occurred while updating the playing message: %s", err.Error())
		}
	}

	log.Infof("[%s:%s] Ready!", session.State.User.Username, session.State.User.ID)
}
