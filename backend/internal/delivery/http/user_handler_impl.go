package http

import (
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
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
	userID := c.Locals("userId").(string)

	userResponse, err := handler.userUsecase.GetCurrentUser(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(model.WebResponse[model.UserResponse]{
		Data: userResponse,
	})
}
