package mock

import (
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/entity"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockRefreshTokenRepository struct {
	Mock mock.Mock
}

func (m *MockRefreshTokenRepository) Save(tx *gorm.DB, refreshToken entity.RefreshToken) error {
	args := m.Mock.Called(tx, refreshToken)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) FindByToken(tx *gorm.DB, token string) (entity.RefreshToken, error) {
	args := m.Mock.Called(tx, token)
	if args.Get(0) == nil {
		return entity.RefreshToken{}, args.Error(1)
	}
	return args.Get(0).(entity.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) DeleteByUserID(tx *gorm.DB, userID string) error {
	args := m.Mock.Called(tx, userID)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteByToken(tx *gorm.DB, token string) error {
	args := m.Mock.Called(tx, token)
	return args.Error(0)
}
