package middleware

import (
	"strings"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/util"
	"github.com/gofiber/fiber/v3"
)

func NewAuthMiddleware(jwtUtil util.JwtUtil) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return fiber.ErrUnauthorized
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwtUtil.ValidateToken(c.Context(), token)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		c.Locals("userId", claims["user_id"].(string))
		c.Locals("role", claims["role"].(string))
		return c.Next()
	}
}
