package http

import (
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandlerImpl struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandlerImpl(userUsecase usecase.UserUsecase) *UserHandlerImpl {
	return &UserHandlerImpl{
		userUsecase: userUsecase,
	}
}

func (handler *UserHandlerImpl) GetCurrent(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(jwt.MapClaims)
	if !ok {
		return fiber.ErrUnauthorized
	}

	userResponse := handler.userUsecase.GetCurrentUser(claims)

	return c.Status(fiber.StatusOK).JSON(model.WebResponse[model.UserResponse]{
		Data: userResponse,
	})
}
