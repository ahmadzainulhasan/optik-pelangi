package repository

import (
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/entity"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	Repository[entity.User]
}

func NewUserRepositoryImpl() *UserRepositoryImpl {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) FindByEmail(tx *gorm.DB, email string) (entity.User, error) {
	var user entity.User
	err := tx.Where("email = ?", email).First(&user).Error
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}
