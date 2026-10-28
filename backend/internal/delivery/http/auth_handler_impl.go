package http

import (
	"errors"
	"strings"
	"time"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
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
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    loginResult.State,
		Expires:  time.Now().Add(5 * time.Minute),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/auth/callback",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "pkce_verifier",
		Value:    loginResult.Verifier,
		Expires:  time.Now().Add(5 * time.Minute),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/auth/callback",
	})

	return c.Redirect().To(loginResult.URL)
}

func (handler *AuthHandlerImpl) GoogleCallback(c fiber.Ctx) error {
	defer func() {
		c.Cookie(&fiber.Cookie{
			Name:     "oauth_state",
			Expires:  fasthttp.CookieExpireDelete,
			HTTPOnly: true,
			Secure:   true,
			SameSite: fiber.CookieSameSiteLaxMode,
			Path:     "/auth/callback",
		})
		c.Cookie(&fiber.Cookie{
			Name:     "pkce_verifier",
			Expires:  fasthttp.CookieExpireDelete,
			HTTPOnly: true,
			Secure:   true,
			SameSite: fiber.CookieSameSiteLaxMode,
			Path:     "/auth/callback",
		})
	}()

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

func (handler *AuthHandlerImpl) Logout(c fiber.Ctx) error {
	defer func() {
		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Expires:  fasthttp.CookieExpireDelete,
			HTTPOnly: true,
			Secure:   true,
			SameSite: fiber.CookieSameSiteStrictMode,
			Path:     "/auth",
		})
	}()

	refreshToken := c.Cookies("refresh_token")
	var accessToken string
	authHeader := c.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		accessToken = strings.TrimPrefix(authHeader, "Bearer ")
	}

	if err := handler.authUsecase.Logout(c.Context(), refreshToken, accessToken); err != nil {
		return err
	}

	return c.JSON(model.WebResponse[string]{
		Data: "OK",
	})
}
