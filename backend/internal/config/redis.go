package config

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func NewRedis(viper *viper.Viper, logger *zerolog.Logger) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("REDIS_ADDR"),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB:       viper.GetInt("REDIS_DB"),
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		logger.Warn().Err(err).Msg("failed to connect redis")
	}

	return redisClient
}
