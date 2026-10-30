package usecase

import (
	"context"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/repository"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type UserUsecaseImpl struct {
	DB             *gorm.DB
	userRepository repository.UserRepository
	log            *zerolog.Logger
}

func NewUserUsecaseImpl(db *gorm.DB, userRepository repository.UserRepository, log *zerolog.Logger) *UserUsecaseImpl {
	return &UserUsecaseImpl{
		DB:             db,
		userRepository: userRepository,
		log:            log,
	}
}

func (usecase *UserUsecaseImpl) GetCurrentUser(ctx context.Context, userID string) (model.UserResponse, error) {
	user, err := usecase.userRepository.FindByID(usecase.DB.WithContext(ctx), userID)
	if err != nil {
		usecase.log.Error().Err(err).Msg("Failed to find user by id")
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

