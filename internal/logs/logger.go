package logs

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

const (
	// MbtLogLevel tool log identifier
	MbtLogLevel = "MBT_LOG_LEVEL"
	// DefLvl default level - should be error
	DefLvl = "info"
)

// Logger expose for usage
var Logger *logrus.Logger

// NewLogger - init logger
func NewLogger() *logrus.Logger {

	var level logrus.Level
	lvl := getLogLevel()
	// In case level doesn't set will not print any message
	level = logLevel(lvl)
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

func logLevel(lvl string) logrus.Level {

	switch lvl {
	case "debug":
		// Used for tracing
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "error":
		return logrus.ErrorLevel
	case "warn":
		return logrus.WarnLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		panic(fmt.Sprintf("the specified log level <%v> is not supported", lvl))
	}
}
