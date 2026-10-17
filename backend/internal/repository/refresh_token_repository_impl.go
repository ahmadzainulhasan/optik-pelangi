package repository

import (
	"time"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/entity"
	"gorm.io/gorm"
)

type RefreshTokenRepositoryImpl struct {
	Repository[entity.RefreshToken]
}

func NewRefreshTokenRepositoryImpl() *RefreshTokenRepositoryImpl {
	return &RefreshTokenRepositoryImpl{}
}

func (repository *RefreshTokenRepositoryImpl) FindByToken(tx *gorm.DB, token string) (entity.RefreshToken, error) {
	var refreshToken entity.RefreshToken
	err := tx.Where("token = ? AND expires_at > ?", token, time.Now()).First(&refreshToken).Error
	if err != nil {
		return entity.RefreshToken{}, err
	}
	return refreshToken, nil
}

func (repository *RefreshTokenRepositoryImpl) DeleteByUserID(tx *gorm.DB, userID string) error {
	return tx.Where("user_id = ?", userID).Delete(&entity.RefreshToken{}).Error
}

func (repository *RefreshTokenRepositoryImpl) DeleteByToken(tx *gorm.DB, token string) error {
	return tx.Where("token = ?", token).Delete(&entity.RefreshToken{}).Error
}
