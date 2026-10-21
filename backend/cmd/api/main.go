package main

import (
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/config"
	"github.com/gofiber/fiber/v3"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger()
	db := config.NewDatabase(viperConfig, log)
	redis := config.NewRedis(viperConfig, log)
	oauth2 := config.NewOAuth2(viperConfig)
	translator := config.NewTranslator()
	validator := config.NewValidator(translator)
	app := config.NewFiber(viperConfig, translator)

	config.Bootstrap(&config.BootstrapConfig{
		DB:          db,
		Redis:       redis,
		OAuth2:      oauth2,
		App:         app,
		Log:         log,
		ViperConfig: viperConfig,
		Validator:   validator,
		Translator:  translator,
	})

	err := app.Listen(":3000", fiber.ListenConfig{
		EnablePrefork: true,
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
