package route

import (
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/delivery/http"
	"github.com/gofiber/fiber/v3"
)

type RouteConfig struct {
	App            *fiber.App
	AuthHandler    http.AuthHandler
	UserHandler    http.UserHandler
	AuthMiddleware fiber.Handler
}

func (config *RouteConfig) Setup() {
	auth := config.App.Group("/auth")
	auth.Get("/login", config.AuthHandler.GoogleLogin)
	auth.Get("/callback", config.AuthHandler.GoogleCallback)
	auth.Post("/refresh", config.AuthHandler.RefreshToken)
	auth.Delete("/logout", config.AuthHandler.Logout)

	user := config.App.Group("/users")
	user.Get("/current", config.AuthMiddleware, config.UserHandler.GetCurrent)
}
