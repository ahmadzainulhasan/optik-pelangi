package repository

import (
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(tx *gorm.DB, email string) (entity.User, error)
	Save(tx *gorm.DB, user entity.User) error
}
