package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
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
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type AuthUsecaseImpl struct {
	DB                     *gorm.DB
	redis                  *redis.Client
	userRepository         repository.UserRepository
	refreshTokenRepository repository.RefreshTokenRepository
	jwtUtil                util.JwtUtil
	validator              *validator.Validate
	log                    *zerolog.Logger
	googleOAuthClient      google.OAuthClient
}

func NewAuthUsecaseImpl(
	db *gorm.DB,
	redis *redis.Client,
	userRepository repository.UserRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	jwtUtil util.JwtUtil,
	validator *validator.Validate,
	log *zerolog.Logger,
	googleOAuthClient google.OAuthClient,
) *AuthUsecaseImpl {
	return &AuthUsecaseImpl{
		DB:                     db,
		redis:                  redis,
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
		return model.GoogleLoginResult{}, fiber.ErrInternalServerError
	}
	state := base64.RawURLEncoding.EncodeToString(stateBytes)

	verifierBytes := make([]byte, 32)
	if _, err := rand.Read(verifierBytes); err != nil {
		usecase.log.Error().Err(err).Msg("Failed to generate pkce verifier")
		return model.GoogleLoginResult{}, fiber.ErrInternalServerError
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

func (usecase *AuthUsecaseImpl) RotateToken(ctx context.Context, refreshToken string) (model.AuthResponse, error) {
	if refreshToken == "" {
		return model.AuthResponse{}, fiber.ErrUnauthorized
	}

	cached, err := usecase.redis.Get(ctx, refreshToken).Bytes()
	if err == nil {
		var authResp model.AuthResponse
		if err := json.Unmarshal(cached, &authResp); err == nil {
			return authResp, nil
		}
	}

	var authResponse model.AuthResponse
	err = usecase.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tokenEntity, err := usecase.refreshTokenRepository.FindByToken(tx, refreshToken)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fiber.ErrUnauthorized
			}
			usecase.log.Error().Err(err).Msg("Failed to find refresh token")
			return err
		}

		if time.Now().After(tokenEntity.ExpiresAt) {
			return fiber.ErrUnauthorized
		}

		user, err := usecase.userRepository.FindByID(tx, tokenEntity.UserID)
		if err != nil {
			usecase.log.Error().Err(err).Msg("Failed to find user for refresh token")
			return err
		}

		authResponse.AccessToken, err = usecase.jwtUtil.GenerateAccessToken(ctx, user)
		if err != nil {
			usecase.log.Error().Err(err).Msg("Failed to generate access token")
			return err
		}

		authResponse.RefreshToken = uuid.New().String()
		tokenEntity.Token = authResponse.RefreshToken
		tokenEntity.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
		if err := usecase.refreshTokenRepository.Save(tx, tokenEntity); err != nil {
			usecase.log.Error().Err(err).Msg("Failed to update refresh token in database")
			return err
		}

		return nil
	})

	if err != nil {
		return model.AuthResponse{}, err
	}

	data, err := json.Marshal(authResponse)
	if err != nil {
		usecase.log.Error().Err(err).Msg("Failed to marshal auth response")
		return model.AuthResponse{}, err
	}

	if err := usecase.redis.Set(ctx, refreshToken, data, 10*time.Second).Err(); err != nil {
		usecase.log.Error().Err(err).Msg("Failed to store refresh token in redis")
		return model.AuthResponse{}, err
	}

	return authResponse, nil
}

func (usecase *AuthUsecaseImpl) Logout(ctx context.Context, refreshToken string, accessToken string) {
	_ = usecase.redis.Del(ctx, accessToken).Err()
	_ = usecase.refreshTokenRepository.DeleteByToken(usecase.DB.WithContext(ctx), refreshToken)
}
