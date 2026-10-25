package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/entity"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/gateway/google"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/repository"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/util"
	"github.com/go-playground/validator/v10"
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

func (usecase *AuthUsecaseImpl) GetGoogleLogin() (model.GoogleLoginResult, error) {
	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		usecase.log.Error().Err(err).Msg("Failed to generate oauth state")
		return model.GoogleLoginResult{}, err
	}
	state := base64.RawURLEncoding.EncodeToString(stateBytes)

	verifierBytes := make([]byte, 32)
	if _, err := rand.Read(verifierBytes); err != nil {
		usecase.log.Error().Err(err).Msg("Failed to generate pkce verifier")
		return model.GoogleLoginResult{}, err
	}
	verifier := base64.RawURLEncoding.EncodeToString(verifierBytes)

	hash := sha256.Sum256([]byte(verifier))

	return model.GoogleLoginResult{
		URL: usecase.googleOAuthClient.GetLoginURL(
			state,
			base64.RawURLEncoding.EncodeToString(hash[:]),
		),
		State:    state,
		Verifier: verifier,
	}, nil
}

func (usecase *AuthUsecaseImpl) ProcessGoogleCallback(ctx context.Context, request model.GoogleCallbackRequest) (model.AuthResponse, error) {
	if err := usecase.validator.Struct(request); err != nil {
		return model.AuthResponse{}, err
	}

	userInfo, err := usecase.googleOAuthClient.GetUserInfo(ctx, request.Code, request.Verifier)
	if err != nil {
		usecase.log.Error().Err(err).Msg("Failed to authenticate with Google")
		return model.AuthResponse{}, err
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

		authResponse.AccessToken, err = usecase.jwtUtil.GenerateAccessToken(ctx, user)
		if err != nil {
			usecase.log.Error().Err(err).Msg("Failed to generate access token")
			return err
		}

		authResponse.RefreshToken = uuid.New().String()
		if err := usecase.refreshTokenRepository.Save(tx, entity.RefreshToken{
			ID:        uuid.New().String(),
			UserID:    user.ID,
			Token:     authResponse.RefreshToken,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		}); err != nil {
			usecase.log.Error().Err(err).Msg("Failed to store refresh token in database")
			return err
		}

		return nil
	})

	if err != nil {
		return model.AuthResponse{}, err
	}

	return authResponse, nil
}
