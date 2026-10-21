package config

import (
	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName(".env")
	config.SetConfigType("env")
	config.AddConfigPath(".")
	config.AddConfigPath("..")
	config.AddConfigPath("../..")

	_ = config.ReadInConfig()

	config.AutomaticEnv()

	return config
}
