package middleware

import "github.com/gofiber/fiber/v3"

func NewAdminMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		if c.Locals("role") != "admin" {
			return fiber.ErrForbidden
		}
		return c.Next()
	}
}
