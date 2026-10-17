package config

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewLogger() *zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	log.Logger = logger

	return &logger
}
