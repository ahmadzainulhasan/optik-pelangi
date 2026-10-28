package usecase

import (
	"context"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
)

type AuthUsecase interface {
	GetGoogleLogin() (model.GoogleLoginResult, error)
	ProcessGoogleCallback(ctx context.Context, request model.GoogleCallbackRequest) (model.AuthResponse, error)
	RotateToken(ctx context.Context, refreshToken string) (model.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string, accessToken string)
}
