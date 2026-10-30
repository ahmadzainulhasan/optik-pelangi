package test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/entity"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestGetCurrentUserSuccess(t *testing.T) {
	mockJwtUtil.Mock.On("ValidateToken", mock.Anything, "valid-bearer-token-in-memory").Return(jwt.MapClaims{
		"user_id": "c3b0f5e1-8f4b-4b2a-9e1a-5f8d2b3c4d5e",
		"role":    "customer",
	}, nil)

	mockUserRepo.Mock.On("FindByID", mock.Anything, "c3b0f5e1-8f4b-4b2a-9e1a-5f8d2b3c4d5e").Return(entity.User{
		ID:      "c3b0f5e1-8f4b-4b2a-9e1a-5f8d2b3c4d5e",
		Name:    "Budi Santoso",
		Email:   "budi.santoso@gmail.com",
		Picture: "https://lh3.googleusercontent.com/a/ACg8ocIq7Q8m_budi_avatar=s96-c",
		Role:    "customer",
	}, nil)

	request := httptest.NewRequest(fiber.MethodGet, "/users/current", nil)
	request.Header.Set("Authorization", "Bearer valid-bearer-token-in-memory")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[model.UserResponse]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)

	assert.Equal(t, "c3b0f5e1-8f4b-4b2a-9e1a-5f8d2b3c4d5e", webResponse.Data.ID)
	assert.Equal(t, "Budi Santoso", webResponse.Data.Name)
	assert.Equal(t, "budi.santoso@gmail.com", webResponse.Data.Email)
	assert.Equal(t, "https://lh3.googleusercontent.com/a/ACg8ocIq7Q8m_budi_avatar=s96-c", webResponse.Data.Picture)
	assert.Equal(t, "customer", webResponse.Data.Role)
}

func TestGetCurrentUserUnauthorizedNoHeader(t *testing.T) {
	request := httptest.NewRequest(fiber.MethodGet, "/users/current", nil)
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
}

func TestGetCurrentUserUnauthorizedInvalidToken(t *testing.T) {
	mockJwtUtil.Mock.On("ValidateToken", mock.Anything, "invalid-or-expired-token").Return(nil, fiber.ErrUnauthorized)

	request := httptest.NewRequest(fiber.MethodGet, "/users/current", nil)
	request.Header.Set("Authorization", "Bearer invalid-or-expired-token")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
}

func TestGetCurrentUserUnauthorizedInvalidHeaderPrefix(t *testing.T) {
	request := httptest.NewRequest(fiber.MethodGet, "/users/current", nil)
	request.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
}

func TestGetCurrentUserUnauthorizedEmptyBearerToken(t *testing.T) {
	mockJwtUtil.Mock.On("ValidateToken", mock.Anything, "").Return(nil, fiber.ErrUnauthorized)

	request := httptest.NewRequest(fiber.MethodGet, "/users/current", nil)
	request.Header.Set("Authorization", "Bearer ")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
}

func TestGetCurrentUserSuccessWithoutPicture(t *testing.T) {
	mockJwtUtil.Mock.On("ValidateToken", mock.Anything, "partial-claims-token").Return(jwt.MapClaims{
		"user_id": "d4c1f6e2-9f5b-4c3a-9e2b-6f9e3c4d5e6f",
		"role":    "customer",
	}, nil)

	mockUserRepo.Mock.On("FindByID", mock.Anything, "d4c1f6e2-9f5b-4c3a-9e2b-6f9e3c4d5e6f").Return(entity.User{
		ID:      "d4c1f6e2-9f5b-4c3a-9e2b-6f9e3c4d5e6f",
		Email:   "user.nopicture@gmail.com",
		Role:    "customer",
	}, nil)

	request := httptest.NewRequest(fiber.MethodGet, "/users/current", nil)
	request.Header.Set("Authorization", "Bearer partial-claims-token")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[model.UserResponse]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)

	assert.Equal(t, "d4c1f6e2-9f5b-4c3a-9e2b-6f9e3c4d5e6f", webResponse.Data.ID)
	assert.Equal(t, "", webResponse.Data.Name)
	assert.Equal(t, "user.nopicture@gmail.com", webResponse.Data.Email)
	assert.Equal(t, "", webResponse.Data.Picture)
	assert.Equal(t, "customer", webResponse.Data.Role)
}

func TestGetCurrentUserNotFound(t *testing.T) {
	mockJwtUtil.Mock.On("ValidateToken", mock.Anything, "user-deleted-token").Return(jwt.MapClaims{
		"user_id": "99999999-9999-9999-9999-999999999999",
		"role":    "customer",
	}, nil)

	mockUserRepo.Mock.On("FindByID", mock.Anything, "99999999-9999-9999-9999-999999999999").Return(entity.User{}, gorm.ErrRecordNotFound)

	request := httptest.NewRequest(fiber.MethodGet, "/users/current", nil)
	request.Header.Set("Authorization", "Bearer user-deleted-token")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
}

func TestGetCurrentUserDatabaseError(t *testing.T) {
	mockJwtUtil.Mock.On("ValidateToken", mock.Anything, "db-error-token").Return(jwt.MapClaims{
		"user_id": "88888888-8888-8888-8888-888888888888",
		"role":    "customer",
	}, nil)

	mockUserRepo.Mock.On("FindByID", mock.Anything, "88888888-8888-8888-8888-888888888888").Return(entity.User{}, errors.New("db error"))

	request := httptest.NewRequest(fiber.MethodGet, "/users/current", nil)
	request.Header.Set("Authorization", "Bearer db-error-token")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
}

