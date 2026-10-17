package http

import (
	"time"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

type AuthHandlerImpl struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandlerImpl(authUsecase usecase.AuthUsecase) *AuthHandlerImpl {
	return &AuthHandlerImpl{
		authUsecase: authUsecase,
	}
}

func (handler *AuthHandlerImpl) GoogleLogin(c fiber.Ctx) error {
	loginURL := handler.authUsecase.GetGoogleLoginURL()
	return c.Redirect().To(loginURL)
}

func (handler *AuthHandlerImpl) GoogleCallback(c fiber.Ctx) error {
	code := c.Query("code")

	authResponse, err := handler.authUsecase.ProcessGoogleCallback(c.Context(), code)
	if err != nil {		
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    authResponse.AccessToken,
		Expires:  time.Now().Add(time.Minute * 60),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    authResponse.RefreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
	})	

	return c.Redirect().To("/")
}
