package config

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(viper *viper.Viper, log *zerolog.Logger) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		viper.GetString("DB_HOST"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_NAME"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_SSLMODE"),
		viper.GetString("DB_TIMEZONE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(&zerologWriter{Logger: log}, logger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect database")
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database instance")
	}

	connection.SetMaxIdleConns(5)
	connection.SetMaxOpenConns(20)
	connection.SetConnMaxLifetime(60 * time.Minute)

	return db
}

type zerologWriter struct {
	Logger *zerolog.Logger
}

func (writer *zerologWriter) Printf(message string, args ...interface{}) {
	writer.Logger.Trace().Msgf(message, args...)
}
