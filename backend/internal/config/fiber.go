package config

import (
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/spf13/viper"
)

func NewErrorHandler(translator ut.Translator) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {

		if fiberError, ok := err.(*fiber.Error); ok {
			return c.Status(fiberError.Code).JSON(model.WebResponse[string]{
				Errors: fiberError.Message,
			})
		}

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			messages := make(map[string]string)
			for _, fieldError := range validationErrors {
				messages[fieldError.Field()] = fieldError.Translate(translator)
			}
			return c.Status(fiber.StatusBadRequest).JSON(model.WebResponse[map[string]string]{
				Data: messages,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{
			Data:   nil,
			Errors: "Terjadi kesalahan pada server, silakan coba lagi nanti",
		})
	}
}

func NewFiber(viper *viper.Viper, translator ut.Translator) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: NewErrorHandler(translator),
	})

	app.Use(recover.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{viper.GetString("CLIENT_ORIGIN")},
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Accept", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
	}))

	return app
}
