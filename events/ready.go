package events

import (
	"discordbot/utils"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

const (
	statusIntervalPeriod = 30 * time.Second
)

var (
	statusTexts = []string{
		"hide and seek",
		"Hunt Royale",
	}
	StatusInterval chan bool
)

func Ready(session *discordgo.Session, event *discordgo.Ready) {
	err := session.UpdateGameStatus(0, utils.RandomString(statusTexts))
	if err != nil {
		log.Warningf("Unable to set status: %s", err.Error())
	}

	StatusInterval = utils.SetInterval(func() {
		err := session.UpdateGameStatus(0, utils.RandomString(statusTexts))
		if err != nil {
			log.Warningf("An error occurred while updating the playing message: %s", err.Error())
		}
	}, statusIntervalPeriod)

	log.Infof("[%s:%s] Ready!", session.State.User.Username, session.State.User.ID)
}
