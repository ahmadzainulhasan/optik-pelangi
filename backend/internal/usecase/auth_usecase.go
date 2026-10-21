package usecase

import (
	"context"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
)

type AuthUsecase interface {
	GetGoogleLoginURL() string
	ProcessGoogleCallback(ctx context.Context, code string) model.AuthResponse
}
