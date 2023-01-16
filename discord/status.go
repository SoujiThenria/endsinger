package discord

import (
	"errors"
	"time"

	log "github.com/SoujiThenria/endsinger/logging"
)

const (
	StatusNone      StatusType = 0
	StatusListen    StatusType = 1
	StatusGame      StatusType = 2
	StatusStreaming StatusType = 3
)

var (
	statusMap = map[string]string{
		"The Global Innovator of Taste": "",
		"Last Resort":                   "",
		"さくら":                           "",
		"Common Sense":                  "",
		"Intoxicated By Youth":          "",
	}
)

// Set the status of the discord bot.
func SetStatus(sType StatusType, Status string) (err error) {
	status.Type = sType
	status.String = Status
	err = status.update()
	return
}

// Update the status of the bot
func (s *Status) update() (err error) {
	// Setting no status
	if s.Type == StatusNone {
		return
	}
	// Cannot set an empty status
	if s.String == "" {
		return errors.New("Cannot update status without a string to display")
	}
	switch s.Type {
	case StatusListen:
		err = session.UpdateListeningStatus(s.String)
	case StatusGame:
		err = session.UpdateGameStatus(1, s.String)
	case StatusStreaming:
		err = session.UpdateStreamingStatus(1, s.String, s.URL)
	}
	return
}

func statusSwitcher(statusType StatusType, stop chan struct{}) {
	ticker := time.NewTicker(1 * time.Hour)
	status.Type = statusType
	for {
		for s, u := range statusMap {
			log.Debug("Changing status to: %s", s)
			status.String = s
			status.URL = u
			err := status.update()
			if err != nil {
				log.Warn("Caanot change status to: %s - Error: %s", s, err)
			}
			select {
			case <-ticker.C:
			case <-stop:
				return
			}
		}
	}
}
