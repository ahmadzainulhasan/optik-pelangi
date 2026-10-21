package util

import (
	"context"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/entity"
	"github.com/golang-jwt/jwt/v5"
)

type JwtUtil interface {
	GenerateAccessToken(ctx context.Context, user entity.User) (string, error)
	ValidateToken(ctx context.Context, accessToken string) (jwt.MapClaims, error)
}
