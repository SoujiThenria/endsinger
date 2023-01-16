package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/SoujiThenria/endsinger/logging"
	"github.com/spf13/viper"
)

type Config struct {
	BotToken   string
	BotStatus  string
	DBPath     string
	ConfigFile string
	PidFile    string
	LogLevel   log.LogLevel
	LogColor   bool
	Daemon     bool
}

// Read flags and the config file.
func getConfig() *Config {
	var err error
	c := new(Config)

	// Set and get flags. None of these options can be set via the config file.
	flag.StringVar(&c.ConfigFile, "c", "/usr/local/endsinger/endsinger.conf", "The path to the config file.")
	flag.StringVar(&c.PidFile, "p", "/var/run/endsinger/endsinger.pid", "The path to the pid file.")
	flag.BoolVar(&c.Daemon, "d", false, "Start the bot as a daemon process.")
	flag.Parse()

	// Resolve relative paths from the flags to absolute ones.
	c.ConfigFile, err = filepath.Abs(c.ConfigFile)
	if err != nil {
		fmt.Println("Cannot convert the config file path:", err)
		os.Exit(1)
	}
	c.PidFile, err = filepath.Abs(c.PidFile)
	if err != nil {
		fmt.Println("Cannot convert the pif file path:", err)
		os.Exit(1)
	}

	// Set the config file parameters.
	viper.SetConfigType("toml")
	viper.SetConfigFile(c.ConfigFile)
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("Cannot read the config file:", err)
		os.Exit(1)
	}

	// Get configuration from the config file.
	c.BotToken = viper.GetString("discord.token")
	c.BotStatus = viper.GetString("discord.status")
	c.DBPath = viper.GetString("discord.database")
	logLevelString := viper.GetString("log.level")
	c.LogColor = viper.GetBool("log.color")

	if c.BotToken == "" {
		fmt.Println("No bot token was found in the config file.")
		os.Exit(1)
	}

	// Convert the log level string to the corresponding log level from the package.
	switch logLevelString {
	case "DEBUG":
		c.LogLevel = log.LevelDebug
	case "INFO":
		c.LogLevel = log.LevelInfo
	case "WARN":
		c.LogLevel = log.LevelWarn
	case "ERROR":
		c.LogLevel = log.LevelError
	case "FATAL":
		c.LogLevel = log.LevelFatal
	default:
		c.LogLevel = log.LevelInfo
	}

	return c
}
