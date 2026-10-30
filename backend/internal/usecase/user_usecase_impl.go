package usecase

import (
	"context"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/repository"
	"gorm.io/gorm"
)

type UserUsecaseImpl struct {
	DB             *gorm.DB
	userRepository repository.UserRepository
}

func NewUserUsecaseImpl(db *gorm.DB, userRepository repository.UserRepository) *UserUsecaseImpl {
	return &UserUsecaseImpl{
		DB:             db,
		userRepository: userRepository,
	}
}

func (usecase *UserUsecaseImpl) GetCurrentUser(ctx context.Context, userID string) (model.UserResponse, error) {
	user, err := usecase.userRepository.FindByID(usecase.DB.WithContext(ctx), userID)
	if err != nil {
		return model.UserResponse{}, err
	}

	return model.UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Picture: user.Picture,
		Role:    user.Role,
	}, nil
}

