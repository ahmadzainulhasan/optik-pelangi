package route

import (
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/delivery/http"
	"github.com/gofiber/fiber/v3"
)

type RouteConfig struct {
	App         *fiber.App
	AuthHandler http.AuthHandler
}

func (config *RouteConfig) Setup() {
	config.SetupAuthRoute()
}

func (config *RouteConfig) SetupAuthRoute() {
	group := config.App.Group("/api/auth")

	group.Get("/login", config.AuthHandler.GoogleLogin)
	group.Get("/callback", config.AuthHandler.GoogleCallback)
}
