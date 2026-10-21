package usecase

import (
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

type UserUsecaseImpl struct {
}

func NewUserUsecaseImpl() *UserUsecaseImpl {
	return &UserUsecaseImpl{}
}

func (usecase *UserUsecaseImpl) GetCurrentUser(claims jwt.MapClaims) model.UserResponse {
	var id, name, email, picture string
	if v, ok := claims["user_id"].(string); ok {
		id = v
	}
	if v, ok := claims["name"].(string); ok {
		name = v
	}
	if v, ok := claims["email"].(string); ok {
		email = v
	}
	if v, ok := claims["picture"].(string); ok {
		picture = v
	}

	return model.UserResponse{
		ID:      id,
		Name:    name,
		Email:   email,
		Picture: picture,
	}
}

