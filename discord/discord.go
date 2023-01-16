package discord

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

var (
	session *discordgo.Session
	status  *Status
	stop    chan struct{}
)

// Initialize the discord package and creates the session
func New(d *Data) (err error) {
	if d.Token == "" {
		return errors.New("The bot token is missing.")
	}
	if d.Status == nil {
		d.Status = &Status{
			Type: StatusNone,
		}
	}

	session, err = discordgo.New("Bot " + d.Token)
	if err != nil {
		return
	}
	status = d.Status
	return
}

// Connect to the Discord gateway and
// add application commands
func Startup() (err error) {
	initHandlers()
	if err = session.Open(); err != nil {
		return
	}
	registerCommands()
	err = status.update()
	stop = make(chan struct{})
	go channelCleanup(stop)
	go statusSwitcher(StatusListen, stop)
	return
}

// Disconnect from the Discord gateway and remove
// all application commands
func Shutdown() (err error) {
	close(stop)
	removeCommands()
	err = session.Close()
	return
}
