package usecase

import (
	"context"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
)

type AuthUsecase interface {
	GetGoogleLogin() (model.GoogleLoginResult, error)
	ProcessGoogleCallback(ctx context.Context, request model.GoogleCallbackRequest) (model.AuthResponse, error)
}
