package http

import "github.com/gofiber/fiber/v3"

type UserHandler interface {
	GetCurrent(c fiber.Ctx) error
}
