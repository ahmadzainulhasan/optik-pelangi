package test

import (
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/config"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/delivery/http"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/delivery/http/middleware"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/delivery/http/route"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/usecase"
	testMock "github.com/ahmadzainulhasan/optik-pelangi/backend/test/mock"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var app *fiber.App
var db *gorm.DB
var viperConfig *viper.Viper
var log *zerolog.Logger
var validate *validator.Validate
var translator ut.Translator
var mockUserRepo *testMock.MockUserRepository
var mockRefreshTokenRepo *testMock.MockRefreshTokenRepository
var mockOAuthClient *testMock.MockGoogleOAuthClient
var mockJwtUtil *testMock.MockJwtUtil

func init() {
	viperConfig = config.NewViper()
	log = config.NewLogger()
	db = config.NewDatabase(viperConfig, log)
	translator = config.NewTranslator()
	validate = config.NewValidator(translator)
	app = config.NewFiber(viperConfig, translator)
	mockUserRepo = new(testMock.MockUserRepository{Mock: mock.Mock{}})
	mockRefreshTokenRepo = new(testMock.MockRefreshTokenRepository{Mock: mock.Mock{}})
	mockOAuthClient = new(testMock.MockGoogleOAuthClient{Mock: mock.Mock{}})
	mockJwtUtil = new(testMock.MockJwtUtil{Mock: mock.Mock{}})

	redisClient := config.NewRedis(viperConfig, log)
	authUsecase := usecase.NewAuthUsecaseImpl(
		db,
		redisClient,
		mockUserRepo,
		mockRefreshTokenRepo,
		mockJwtUtil,
		validate,
		log,
		mockOAuthClient,
	)
	authHandler := http.NewAuthHandlerImpl(authUsecase, viperConfig)

	userUsecase := usecase.NewUserUsecaseImpl(db, mockUserRepo)
	userHandler := http.NewUserHandlerImpl(userUsecase)
	authMiddleware := middleware.NewAuthMiddleware(mockJwtUtil)
	adminMiddleware := middleware.NewAdminMiddleware()

	(&route.RouteConfig{
		App:             app,
		AuthHandler:     authHandler,
		UserHandler:     userHandler,
		AuthMiddleware:  authMiddleware,
		AdminMiddleware: adminMiddleware,
	}).Setup()
}
