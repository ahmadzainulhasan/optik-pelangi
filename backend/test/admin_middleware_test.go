package test

import (
	"net/http/httptest"
	"testing"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/delivery/http/middleware"
	testMock "github.com/ahmadzainulhasan/optik-pelangi/backend/test/mock"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAdminMiddlewareSuccess(t *testing.T) {
	testApp := fiber.New()
	testApp.Get("/admin/test", func(c fiber.Ctx) error {
		c.Locals("role", "admin")
		return c.Next()
	}, middleware.NewAdminMiddleware(), func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	request := httptest.NewRequest(fiber.MethodGet, "/admin/test", nil)
	response, err := testApp.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func TestAdminMiddlewareForbiddenWhenRoleIsCustomer(t *testing.T) {
	testApp := fiber.New()
	testApp.Get("/admin/test", func(c fiber.Ctx) error {
		c.Locals("role", "customer")
		return c.Next()
	}, middleware.NewAdminMiddleware(), func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	request := httptest.NewRequest(fiber.MethodGet, "/admin/test", nil)
	response, err := testApp.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, response.StatusCode)
}

func TestAdminMiddlewareForbiddenWhenNoRoleInLocals(t *testing.T) {
	testApp := fiber.New()
	testApp.Get("/admin/test", middleware.NewAdminMiddleware(), func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	request := httptest.NewRequest(fiber.MethodGet, "/admin/test", nil)
	response, err := testApp.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, response.StatusCode)
}

func TestAdminMiddlewareWithAuthMiddlewareSuccess(t *testing.T) {
	testApp := fiber.New()
	mockJwt := new(testMock.MockJwtUtil{Mock: mock.Mock{}})
	mockJwt.Mock.On("ValidateToken", mock.Anything, "admin-token").Return(jwt.MapClaims{
		"user_id": "c3b0f5e1-8f4b-4b2a-9e1a-5f8d2b3c4d5e",
		"role":    "admin",
	}, nil)

	testApp.Get("/admin/test", middleware.NewAuthMiddleware(mockJwt), middleware.NewAdminMiddleware(), func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	request := httptest.NewRequest(fiber.MethodGet, "/admin/test", nil)
	request.Header.Set("Authorization", "Bearer admin-token")
	response, err := testApp.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func TestAdminMiddlewareWithAuthMiddlewareCustomerForbidden(t *testing.T) {
	testApp := fiber.New()
	mockJwt := new(testMock.MockJwtUtil{Mock: mock.Mock{}})
	mockJwt.Mock.On("ValidateToken", mock.Anything, "customer-token").Return(jwt.MapClaims{
		"user_id": "c3b0f5e1-8f4b-4b2a-9e1a-5f8d2b3c4d5e",
		"role":    "customer",
	}, nil)

	testApp.Get("/admin/test", middleware.NewAuthMiddleware(mockJwt), middleware.NewAdminMiddleware(), func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	request := httptest.NewRequest(fiber.MethodGet, "/admin/test", nil)
	request.Header.Set("Authorization", "Bearer customer-token")
	response, err := testApp.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, response.StatusCode)
}

