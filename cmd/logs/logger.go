package logs

import (
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
	"os"
)

const (
	MbtLogLevel = "MBT_LOG_LEVEL"
	DefLvl      = "info"
)

// Logger - logrus variable
var Logger *logrus.Logger

// NewLogger - init logger
func NewLogger() *logrus.Logger {

	var level logrus.Level
	lvl := getLogLevel()
	// In case level doesn't set will not print any message
	level = logLevel(lvl, level)
	logger := &logrus.Logger{
		Out:   os.Stdout,
		Level: level,
		Formatter: &prefixed.TextFormatter{
			DisableColors:   false,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		},
	}
	Logger = logger
	return Logger
}

// GetLogLevel - Get level from env
func getLogLevel() string {
	// TODO Check env if coming from external config or local
	lvl, _ := os.LookupEnv(MbtLogLevel)
	if lvl != "" {
		return lvl
	}
	return DefLvl
}

func logLevel(lvl string, level logrus.Level) logrus.Level {

	switch lvl {
	case "debug":
		// Used for tracing
		level = logrus.DebugLevel
	case "info":
		level = logrus.InfoLevel
	case "error":
		level = logrus.ErrorLevel
	case "warn":
		level = logrus.WarnLevel
	case "fatal":
		level = logrus.FatalLevel
	case "panic":
		level = logrus.PanicLevel
	default:
		panic("The specified log level is not supported.")
	}
	return level
}
