package http

import (
	"time"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
)

type AuthHandlerImpl struct {
	authUsecase usecase.AuthUsecase
	viper       *viper.Viper
}

func NewAuthHandlerImpl(authUsecase usecase.AuthUsecase, viper *viper.Viper) *AuthHandlerImpl {
	return &AuthHandlerImpl{
		authUsecase: authUsecase,
		viper:       viper,
	}
}

func (handler *AuthHandlerImpl) GoogleLogin(c fiber.Ctx) error {
	loginURL := handler.authUsecase.GetGoogleLoginURL()
	return c.Redirect().To(loginURL)
}

func (handler *AuthHandlerImpl) GoogleCallback(c fiber.Ctx) error {
	code := c.Query("code")

	authResponse := handler.authUsecase.ProcessGoogleCallback(c.Context(), code)
	if authResponse.AccessToken != "" {

		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    authResponse.AccessToken,
			Expires:  time.Now().Add(time.Minute * 60),
			HTTPOnly: true,
			Secure:   true,
			SameSite: fiber.CookieSameSiteStrictMode,
			Path:     "/",
		})

		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    authResponse.RefreshToken,
			Expires:  time.Now().Add(time.Hour * 24 * 7),
			HTTPOnly: true,
			Secure:   true,
			SameSite: fiber.CookieSameSiteStrictMode,
			Path:     "/",
		})
	}

	return c.Redirect().To(handler.viper.GetString("CLIENT_ORIGIN") + "/callback")
}
