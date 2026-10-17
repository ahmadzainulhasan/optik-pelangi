package http

import "github.com/gofiber/fiber/v3"

type AuthHandler interface {
	GoogleLogin(c fiber.Ctx) error
	GoogleCallback(c fiber.Ctx) error
}
