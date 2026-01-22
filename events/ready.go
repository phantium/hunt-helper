package events

import (
	"discordbot/utils"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

const (
	statusIntervalPeriod = 5 * time.Minute // Reduced interval for more frequent updates
)

var (
	statusTexts = []string{
		"Hunt Royale",
		"Dark Forest",
		"Dragon's Dungeon",
		"Kraken's Ship",
		"Yeti's Tundra",
		"Maze",
		"Abyssal Maze",
		"Chaos Dungeon",
	}
)

func Ready(session *discordgo.Session, event *discordgo.Ready) {
	updateStatus := func() {
		status := utils.RandomString(statusTexts)
		err := session.UpdateGameStatus(0, status)
		if err != nil {
			log.Warningf("Unable to set status: %s", err.Error())
		} else {
			log.Infof("Updated game status to: %s", status)
		}
	}

	// Initial status update
	updateStatus()

	// Start a goroutine for periodic updates
	go func() {
		ticker := time.NewTicker(statusIntervalPeriod)
		for range ticker.C {
			updateStatus()
		}
	}()

	log.Infof("[%s:%s] Ready!", session.State.User.Username, session.State.User.ID)
}
