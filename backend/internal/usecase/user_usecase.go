package usecase

import (
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

type UserUsecase interface {
	GetCurrentUser(claims jwt.MapClaims) model.UserResponse
}

