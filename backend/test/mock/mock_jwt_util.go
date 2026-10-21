package mock

import (
	"context"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type MockJwtUtil struct {
	Mock mock.Mock
}

func (m *MockJwtUtil) GenerateAccessToken(ctx context.Context, user entity.User) (string, error) {
	args := m.Mock.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockJwtUtil) ValidateToken(ctx context.Context, tokenString string) (jwt.MapClaims, error) {
	args := m.Mock.Called(ctx, tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(jwt.MapClaims), args.Error(1)
}
