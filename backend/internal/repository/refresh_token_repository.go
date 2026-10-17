package repository

import (
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/entity"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Save(tx *gorm.DB, refreshToken entity.RefreshToken) error
	FindByToken(tx *gorm.DB, token string) (entity.RefreshToken, error)
	DeleteByUserID(tx *gorm.DB, userID string) error
	DeleteByToken(tx *gorm.DB, token string) error
}
