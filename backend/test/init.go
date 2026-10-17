package test

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/config"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/delivery/http"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/delivery/http/route"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/usecase"
	testMock "github.com/ahmadzainulhasan/optik-pelangi/backend/test/mock"
	ut "github.com/go-playground/universal-translator"
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
	var err error
	sqlDB, sqlMock, err = sqlmock.New()
	if err != nil {
		panic(err)
	}

	db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	viperConfig = config.NewViper()
	viperConfig.Set("JWT_SECRET", "super-secret-jwt-key-for-testing-purpose-123456")
	log = config.NewLogger()
	translator = config.NewTranslator()
	validate = config.NewValidator(translator)
	app = config.NewFiber(viperConfig, translator)
	mockUserRepo = new(testMock.MockUserRepository{Mock: mock.Mock{}})
	mockRefreshTokenRepo = new(testMock.MockRefreshTokenRepository{Mock: mock.Mock{}})
	mockOAuthClient = new(testMock.MockGoogleOAuthClient{Mock: mock.Mock{}})
	mockJwtUtil = new(testMock.MockJwtUtil{Mock: mock.Mock{}})

	(&route.RouteConfig{
		App: app,
		AuthHandler: http.NewAuthHandlerImpl(usecase.NewAuthUsecaseImpl(
			db,
			mockUserRepo,
			mockRefreshTokenRepo,
			mockJwtUtil,
			validate,
			log,
			mockOAuthClient,
		)),
	}).Setup()
}
