package logger

import (
	"SubscriptionService/configs"
	"os"

	"github.com/rs/zerolog"
)

func NewLogger(config *configs.LogConfig) *zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.Level(config.Level))
	var logger zerolog.Logger

	if config.Format == "json" {
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	} else {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
		logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
	}

	return &logger
}
