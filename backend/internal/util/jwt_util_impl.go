package util

import (
	"context"
	"time"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/entity"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type JwtUtilImpl struct {
	redis *redis.Client
	viper *viper.Viper
}

func NewJwtUtilImpl(redis *redis.Client, viper *viper.Viper) *JwtUtilImpl {
	return &JwtUtilImpl{
		redis: redis,
		viper: viper,
	}
}

func (util *JwtUtilImpl) GenerateAccessToken(ctx context.Context, user entity.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"picture": user.Picture,
		"exp":     time.Now().Add(15 * time.Minute).UnixMilli(),
	})

	accessToken, err := token.SignedString([]byte(util.viper.GetString("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	if err := util.redis.Set(ctx, accessToken, user.ID, 15*time.Minute).Err(); err != nil {
		return "", err
	}

	return accessToken, nil
}

func (util *JwtUtilImpl) ValidateToken(ctx context.Context, accessToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		return []byte(util.viper.GetString("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	claims := token.Claims.(jwt.MapClaims)
	if int64(claims["exp"].(float64)) < time.Now().UnixMilli() {
		return nil, fiber.ErrUnauthorized
	}

	result, err := util.redis.Exists(ctx, accessToken).Result()
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	if result == 0 {
		return nil, fiber.ErrUnauthorized
	}

	return claims, nil
}
