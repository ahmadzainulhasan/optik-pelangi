package test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCurrentUserSuccessWithCookie(t *testing.T) {
	mockJwtUtil.Mock.On("ValidateToken", mock.Anything, "valid-access-token-cookie").Return(jwt.MapClaims{
		"user_id": "c3b0f5e1-8f4b-4b2a-9e1a-5f8d2b3c4d5e",
		"email":   "budi.santoso@gmail.com",
		"name":    "Budi Santoso",
		"picture": "https://lh3.googleusercontent.com/a/ACg8ocIq7Q8m_budi_avatar=s96-c",
	}, nil)

	request := httptest.NewRequest(fiber.MethodGet, "/users/current", nil)
	request.Header.Set("Cookie", "access_token=valid-access-token-cookie")
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
}

func TestGetCurrentUserUnauthorizedNoCookie(t *testing.T) {
	request := httptest.NewRequest(fiber.MethodGet, "/users/current", nil)
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
}

func TestGetCurrentUserUnauthorizedInvalidToken(t *testing.T) {
	mockJwtUtil.Mock.On("ValidateToken", mock.Anything, "invalid-or-expired-token").Return(nil, fiber.ErrUnauthorized)

	request := httptest.NewRequest(fiber.MethodGet, "/users/current", nil)
	request.Header.Set("Cookie", "access_token=invalid-or-expired-token")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
}

