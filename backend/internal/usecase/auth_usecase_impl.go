package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/entity"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/gateway/google"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/repository"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type AuthUsecaseImpl struct {
	DB                     *gorm.DB
	userRepository         repository.UserRepository
	refreshTokenRepository repository.RefreshTokenRepository
	jwtUtil                util.JwtUtil
	validator              *validator.Validate
	log                    *zerolog.Logger	
	googleOAuthClient      google.OAuthClient
}

func NewAuthUsecaseImpl(
	db *gorm.DB,
	userRepository repository.UserRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	jwtUtil util.JwtUtil,
	validator *validator.Validate,
	log *zerolog.Logger,	
	googleOAuthClient google.OAuthClient,
) *AuthUsecaseImpl {
	return &AuthUsecaseImpl{
		DB:                     db,
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
		jwtUtil:                jwtUtil,
		validator:              validator,
		log:                    log,		
		googleOAuthClient:      googleOAuthClient,
	}
}

func (usecase *AuthUsecaseImpl) GetGoogleLoginURL() string {
	return usecase.googleOAuthClient.GetLoginURL()
}

func (usecase *AuthUsecaseImpl) ProcessGoogleCallback(ctx context.Context, code string) (model.AuthResponse, error) {
	userInfo, err := usecase.googleOAuthClient.GetUserInfo(ctx, code)
	if err != nil {
		usecase.log.Error().Err(err).Msg("Failed to authenticate with Google")
		return model.AuthResponse{}, fiber.NewError(fiber.StatusBadRequest, "Autentikasi Google gagal atau kode otorisasi tidak valid")
	}
	
	var authResponse model.AuthResponse
	err = usecase.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		user, err := usecase.userRepository.FindByEmail(tx, userInfo.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				user = entity.User{
					ID:       uuid.New().String(),
					Email:    userInfo.Email,
					Name:     userInfo.Name,
					GoogleID: userInfo.ID,
					Picture:  userInfo.Picture,
				}
				if err := usecase.userRepository.Save(tx, user); err != nil {
					usecase.log.Error().Err(err).Msg("Failed to create new user")
					return err
				}
			} else {
				usecase.log.Error().Err(err).Msg("Database error during find user")
				return err
			}
		} else {
			user.Name = userInfo.Name
			user.Picture = userInfo.Picture
			user.GoogleID = userInfo.ID
			if err := usecase.userRepository.Save(tx, user); err != nil {
				usecase.log.Error().Err(err).Msg("Failed to update user")
				return err
			}
		}

		accessToken, err := usecase.jwtUtil.GenerateAccessToken(ctx, user.ID, user.Email)
		if err != nil {
			usecase.log.Error().Err(err).Msg("Failed to generate access token")
			return err
		}

		refreshToken := entity.RefreshToken{
			ID:        uuid.New().String(),
			UserID:    user.ID,
			Token:     uuid.New().String(),
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		}

		if err := usecase.refreshTokenRepository.Save(tx, refreshToken); err != nil {
			usecase.log.Error().Err(err).Msg("Failed to store refresh token in database")
			return err
		}

		authResponse = model.AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken.Token,
		}

		return nil
	})

	if err != nil {
		return model.AuthResponse{}, err
	}

	return authResponse, nil
}
