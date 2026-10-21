package config

import (
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/delivery/http"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/delivery/http/route"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/gateway/google"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/repository"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/usecase"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/util"
)

type BootstrapConfig struct {
	DB           *gorm.DB
	Redis        *redis.Client
	OAuth2 		 *oauth2.Config
	App          *fiber.App
	Log          *zerolog.Logger
	ViperConfig  *viper.Viper
	Validator    *validator.Validate
	Translator   ut.Translator
}

func Bootstrap(config *BootstrapConfig) {

	userRepository := repository.NewUserRepositoryImpl()
	refreshTokenRepository := repository.NewRefreshTokenRepositoryImpl()
	
	jwtUtil := util.NewJwtUtilImpl(config.Redis, config.ViperConfig)
	googleOAuthClient := google.NewOAuthClientImpl(config.OAuth2)

	authUsecase := usecase.NewAuthUsecaseImpl(
		config.DB,
		userRepository,
		refreshTokenRepository,
		jwtUtil,
		config.Validator,
		config.Log,		
		googleOAuthClient,
	)
	
	authHandler := http.NewAuthHandlerImpl(authUsecase, config.ViperConfig)

	routeConfig := route.RouteConfig{
		App:         config.App,
		AuthHandler: authHandler,
	}
	routeConfig.Setup()
}
