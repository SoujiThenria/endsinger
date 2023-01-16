package logging

import (
	"io"
	logging "log"
	"os"
)

type LogLevel int8

const (
	// Available log levels
	LevelDebug LogLevel = 1
	LevelInfo  LogLevel = 2
	LevelWarn  LogLevel = 3
	LevelError LogLevel = 4
	LevelFatal LogLevel = 5
)

var (
	log      *logging.Logger
	loglevel LogLevel
	color    bool

	levelString = [...]string{
		"",
		"[DEBUG] ",
		"[INFO] ",
		"[WARN] ",
		"[ERROR] ",
		"[FATAL] ",
	}

	levelStringColor = [...]string{
		"",
		"\x1b[36m[DEBUG]\x1b[0m ", // cyan
		"\x1b[34m[INFO]\x1b[0m ",  // blue
		"\x1b[33m[WARN]\x1b[0m ",  // yellow
		"\x1b[31m[ERROR]\x1b[0m ", // red
		"\x1b[31m[FATAL]\x1b[0m ", // red
	}
)

// Create a new global logger.
func Init(level LogLevel, useColor bool) {
	loglevel = level
	log = logging.New(os.Stderr, "", logging.Ldate|logging.Ltime)
	color = useColor
}

// Set the log level
func SetLogLevel(level LogLevel) {
	loglevel = level
}

// Use color for the log messages.
func UseColor(useColor bool) {
	color = useColor
}

// Set the output io.Writer
func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

// Print logs with timestamp or without.
func Timestamp(ts bool) {
	if ts {
		log.SetFlags(logging.Ldate | logging.Ltime)
		return
	}
	log.SetFlags(0)
}

// General log function
func logf(level LogLevel, format string, args ...interface{}) {
	if level < loglevel {
		return
	}
	if color {
		log.Printf(levelStringColor[level]+format+"\n", args...)
		return
	}

	log.Printf(levelString[level]+format+"\n", args...)
}

func Debug(format string, args ...interface{}) {
	logf(LevelDebug, format, args...)
}

func Info(format string, args ...interface{}) {
	logf(LevelInfo, format, args...)
}

func Warn(format string, args ...interface{}) {
	logf(LevelWarn, format, args...)
}

func Error(format string, args ...interface{}) {
	logf(LevelError, format, args...)
}

func Fatal(format string, args ...interface{}) {
	logf(LevelFatal, format, args...)
	os.Exit(1)
}
