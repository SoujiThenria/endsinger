package main

import (
	"fmt"
	"log/syslog"
	"os"
	"os/signal"
	"syscall"

	"github.com/SoujiThenria/endsinger/database"
	"github.com/SoujiThenria/endsinger/discord"
	log "github.com/SoujiThenria/endsinger/logging"
	"github.com/sevlyar/go-daemon"
)

const (
	VERSION = "1.0.0"
)

func main() {
	// Just for simplicity reasons and consistency.
	var err error

	// Get the configuration for the bot.
	conf := getConfig()

	// The log config changes if the bot is started as daemon.
	log.Init(conf.LogLevel, conf.LogColor)

	// Start as daemon.
	if conf.Daemon {
		// Use the systemlog facility to log to the daemon log file.
		slog, err := syslog.New(syslog.LOG_DAEMON, "endsinger")
		if err != nil {
			fmt.Println("Cannot change the log output to the system log facility.", err)
			os.Exit(1)
		}
		// Disable color and timestamps.
		log.UseColor(false)
		log.Timestamp(false)
		log.SetOutput(slog)

		// Daemon config
		cntx := &daemon.Context{
			PidFileName: conf.PidFile,
			PidFilePerm: 0644,
			Umask:       027,
		}
		d, err := cntx.Reborn()
		if err != nil {
			fmt.Println("Cannot start process as a daemon:", err)
			os.Exit(1)
		}
		if d != nil {
			return
		}
		defer cntx.Release()
		log.Debug("Bot started as a daemon.")
	}

	log.Info("Starting endsinger VERSION-%s", VERSION)

	// Database stuff
	if err = database.New(conf.DBPath); err != nil {
		log.Fatal("Cannot initialize the database: %s", err)
	}
	defer func() {
		if err = database.Close(); err != nil {
			log.Warn("While closing the database: %s", err)
		}
	}()

	// Discord stuff
	err = discord.New(&discord.Data{
		Token: conf.BotToken,
	})

	// Set bot status if there is any.
	if conf.BotStatus != "" {
		err = discord.SetStatus(discord.StatusListen, conf.BotStatus)
		if err != nil {
			log.Warn("Cannot set the bot status: %s", err)
		}
	}
	if err != nil {
		log.Fatal("Cannot initialize discord: %s", err)
	}
	if err = discord.Startup(); err != nil {
		log.Fatal("Cannot start the discord bot: %s", err)
	}
	defer func() {
		if err = discord.Shutdown(); err != nil {
			log.Error("Failed to stop the bot properly: %s", err)
		}
	}()

	// Signal handler
	log.Info("endsinger started.")
	if conf.Daemon {
		err := daemon.ServeSignals()
		if err != nil {
			log.Fatal("Daemon cannot serve signals: %s", err)
		}
	} else {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
	}
	log.Info("Signal received, shutting down.")
}
