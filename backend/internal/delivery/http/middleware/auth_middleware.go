package middleware

import (
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/util"
	"github.com/gofiber/fiber/v3"
)

func NewAuthMiddleware(jwtUtil util.JwtUtil) fiber.Handler {
	return func(c fiber.Ctx) error {
		accessToken := c.Cookies("access_token")
		if accessToken == "" {
			return fiber.ErrUnauthorized
		}

		claims, err := jwtUtil.ValidateToken(c.Context(), accessToken)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		c.Locals("claims", claims)
		return c.Next()
	}
}
