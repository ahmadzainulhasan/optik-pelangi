package test

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var app *fiber.App
var db *gorm.DB
var sqlMock sqlmock.Sqlmock
var viperConfig *viper.Viper
var log *zerolog.Logger
var validate *validator.Validate
var translator ut.Translator
var mockUserRepo *testMock.MockUserRepository
var mockRefreshTokenRepo *testMock.MockRefreshTokenRepository
var mockOAuthClient *testMock.MockGoogleOAuthClient
var mockJwtUtil *testMock.MockJwtUtil

func init() {
	var sqlDB *sql.DB
	sqlDB, sqlMock, _ = sqlmock.New()

	db, _ = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	viperConfig = config.NewViper()
	log = config.NewLogger()
	translator = config.NewTranslator()
	validate = config.NewValidator(translator)
	app = config.NewFiber(viperConfig, translator)
	mockUserRepo = new(testMock.MockUserRepository{Mock: mock.Mock{}})
	mockRefreshTokenRepo = new(testMock.MockRefreshTokenRepository{Mock: mock.Mock{}})
	mockOAuthClient = new(testMock.MockGoogleOAuthClient{Mock: mock.Mock{}})
	mockJwtUtil = new(testMock.MockJwtUtil{Mock: mock.Mock{}})

	authUsecase := usecase.NewAuthUsecaseImpl(
		db,
		mockUserRepo,
		mockRefreshTokenRepo,
		mockJwtUtil,
		validate,
		log,
		mockOAuthClient,
	)
	authHandler := http.NewAuthHandlerImpl(authUsecase, viperConfig)

	userUsecase := usecase.NewUserUsecaseImpl()
	userHandler := http.NewUserHandlerImpl(userUsecase)
	authMiddleware := middleware.NewAuthMiddleware(mockJwtUtil)

	(&route.RouteConfig{
		App:            app,
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		AuthMiddleware: authMiddleware,
	}).Setup()
}
