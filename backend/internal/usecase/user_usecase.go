package usecase

import (
	"context"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
)

type UserUsecase interface {
	GetCurrentUser(ctx context.Context, userID string) (model.UserResponse, error)
}

