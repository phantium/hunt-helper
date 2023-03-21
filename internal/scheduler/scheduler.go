package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
)

var s *gocron.Scheduler

func init() {
	s = gocron.NewScheduler(time.UTC)
	s.StartAsync()
}

func UpdatePostEveryMinute() {
	s.Every(60).Seconds().Do(func() {

	})
}
