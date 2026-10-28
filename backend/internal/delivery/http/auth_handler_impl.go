package http

import (
	"errors"
	"time"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
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
	loginResult, err := handler.authUsecase.GetGoogleLogin()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    loginResult.State,
		Expires:  time.Now().Add(5 * time.Minute),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "pkce_verifier",
		Value:    loginResult.Verifier,
		Expires:  time.Now().Add(5 * time.Minute),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
	})

	return c.Redirect().To(loginResult.URL)
}

func (handler *AuthHandlerImpl) GoogleCallback(c fiber.Ctx) error {
	defer c.ClearCookie("oauth_state", "pkce_verifier")

	var request model.GoogleCallbackRequest
	if err := c.Bind().All(&request); err != nil {
		return fiber.ErrBadRequest
	}

	authResponse, err := handler.authUsecase.ProcessGoogleCallback(c.Context(), request)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    authResponse.RefreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     "/auth",
	})

	return c.JSON(model.WebResponse[model.LoginResponse]{
		Data: model.LoginResponse{
			AccessToken: authResponse.AccessToken,
		},
	})
}

func (handler *AuthHandlerImpl) RefreshToken(c fiber.Ctx) error {
	authResponse, err := handler.authUsecase.RotateToken(c.Context(), c.Cookies("refresh_token"))
	if err != nil {
		if errors.Is(err, fiber.ErrUnauthorized) {
			c.ClearCookie("refresh_token")
		}
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    authResponse.RefreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     "/auth",
	})

	return c.JSON(model.WebResponse[model.LoginResponse]{
		Data: model.LoginResponse{
			AccessToken: authResponse.AccessToken,
		},
	})
}
