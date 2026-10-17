package util

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type JwtUtil interface {
	GenerateAccessToken(ctx context.Context, userID string, email string) (string, error)
	ValidateToken(ctx context.Context, accessToken string) (jwt.MapClaims, error)
}

