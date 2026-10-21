package test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/entity"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestGoogleLogin(t *testing.T) {
	mockOAuthClient.Mock.On("GetLoginURL").Return("https://accounts.google.com/o/oauth2/v2/auth?access_type=offline&client_id=718293847561-abcd1234efgh5678.apps.googleusercontent.com&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fapi%2Fauth%2Fcallback&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.profile&state=state-token").Once()

	request := httptest.NewRequest(fiber.MethodGet, "/api/auth/login", nil)
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusSeeOther, response.StatusCode)
	assert.Equal(t, "https://accounts.google.com/o/oauth2/v2/auth?access_type=offline&client_id=718293847561-abcd1234efgh5678.apps.googleusercontent.com&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fapi%2Fauth%2Fcallback&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.profile&state=state-token", response.Header.Get("Location"))
}

func TestGoogleCallbackSuccessNewUser(t *testing.T) {
	userInfo := model.GoogleUserInfo{
		ID:      "108234918237491823741",
		Email:   "budi.santoso@gmail.com",
		Name:    "Budi Santoso",
		Picture: "https://lh3.googleusercontent.com/a/ACg8ocIq7Q8m_budi_avatar=s96-c",
	}

	mockOAuthClient.Mock.On("GetUserInfo", mock.Anything, "4/0AWgavTe-real_code_new_user_xyz123").Return(userInfo, nil).Once()
	mockUserRepo.Mock.On("FindByEmail", mock.Anything, userInfo.Email).Return(entity.User{}, gorm.ErrRecordNotFound).Once()
	mockUserRepo.Mock.On("Save", mock.Anything, mock.Anything).Return(nil).Once()
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, mock.Anything).Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJiNWYwZDNhMS1kZWIzLTQzNmQtOTYyMS02YzM1NWI3MjMwMWEiLCJlbWFpbCI6ImJ1ZGkuc2FudG9zb0BnbWFpbC5jb20iLCJleHAiOjE3NjAxNDQwMDB9.sig-budi-access-token", nil).Once()
	mockRefreshTokenRepo.Mock.On("Save", mock.Anything, mock.Anything).Return(nil).Once()

	request := httptest.NewRequest(fiber.MethodGet, "/api/auth/callback?code=4/0AWgavTe-real_code_new_user_xyz123", nil)
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, 303, response.StatusCode)
	assert.Equal(t, "http://localhost:5173/auth/callback", response.Header.Get("Location"))
	assert.Len(t, response.Cookies(), 2)
	assert.Equal(t, http.SameSiteLaxMode, response.Cookies()[0].SameSite)
}

func TestGoogleCallbackSuccessExistingUser(t *testing.T) {
	userInfo := model.GoogleUserInfo{
		ID:      "117283948192384719283",
		Email:   "siti.aminah@gmail.com",
		Name:    "Siti Aminah",
		Picture: "https://lh3.googleusercontent.com/a/ACg8ocK82Hs_siti_new_avatar=s96-c",
	}
	user := entity.User{
		ID:       "c3b0f5e1-8f4b-4b2a-9e1a-5f8d2b3c4d5e",
		Email:    userInfo.Email,
		Name:     userInfo.Name,
		GoogleID: userInfo.ID,
		Picture:  userInfo.Picture,
	}

	mockOAuthClient.Mock.On("GetUserInfo", mock.Anything, "4/0AWgavTf-real_code_existing_user_abc456").Return(userInfo, nil).Once()
	mockUserRepo.Mock.On("FindByEmail", mock.Anything, userInfo.Email).Return(user, nil).Once()
	mockUserRepo.Mock.On("Save", mock.Anything, user).Return(nil).Once()
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, user).Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJjM2IwZjVlMS04ZjRiLTRiMmEtOWUxYS01ZjhkMmIzYzRkNWUiLCJlbWFpbCI6InNpdGkuYW1pbmFoQGdtYWlsLmNvbSIsImV4cCI6MTc2MDE0NDAwMH0.sig-siti-access-token", nil).Once()
	mockRefreshTokenRepo.Mock.On("Save", mock.Anything, mock.Anything).Return(nil).Once()

	request := httptest.NewRequest(fiber.MethodGet, "/api/auth/callback?code=4/0AWgavTf-real_code_existing_user_abc456", nil)
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, 303, response.StatusCode)
	assert.Len(t, response.Cookies(), 2)
}

func TestGoogleCallbackMissingCode(t *testing.T) {
	mockOAuthClient.Mock.On("GetUserInfo", mock.Anything, "").Return(model.GoogleUserInfo{}, errors.New("empty authorization code")).Once()

	request := httptest.NewRequest(fiber.MethodGet, "/api/auth/callback", nil)
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
}

func TestGoogleCallbackEmptyCode(t *testing.T) {
	mockOAuthClient.Mock.On("GetUserInfo", mock.Anything, "").Return(model.GoogleUserInfo{}, errors.New("empty authorization code")).Once()

	request := httptest.NewRequest(fiber.MethodGet, "/api/auth/callback?code=", nil)
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
}

func TestGoogleCallbackOAuthExchangeFailed(t *testing.T) {
	mockOAuthClient.Mock.On("GetUserInfo", mock.Anything, "4/0AWgavTg-invalid-or-expired-auth-token-999").
		Return(model.GoogleUserInfo{}, errors.New("oauth2: cannot fetch token: 400 Bad Request\nResponse: {\"error\": \"invalid_grant\", \"error_description\": \"Bad Request\"}")).Once()

	request := httptest.NewRequest(fiber.MethodGet, "/api/auth/callback?code=4/0AWgavTg-invalid-or-expired-auth-token-999", nil)
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
}

func TestGoogleCallbackDatabaseErrorFindUser(t *testing.T) {
	userInfo := model.GoogleUserInfo{
		ID:      "109823749182374918234",
		Email:   "ahmad.dahlan@gmail.com",
		Name:    "Ahmad Dahlan",
		Picture: "https://lh3.googleusercontent.com/a/ACg8ocM93Lk_ahmad_avatar=s96-c",
	}

	mockOAuthClient.Mock.On("GetUserInfo", mock.Anything, "4/0AWgavTh-real_code_find_err_789").Return(userInfo, nil).Once()
	mockUserRepo.Mock.On("FindByEmail", mock.Anything, userInfo.Email).Return(entity.User{}, errors.New("dial tcp 127.0.0.1:3306: connect: connection refused")).Once()

	request := httptest.NewRequest(fiber.MethodGet, "/api/auth/callback?code=4/0AWgavTh-real_code_find_err_789", nil)
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
}

func TestGoogleCallbackDatabaseErrorSaveNewUser(t *testing.T) {
	userInfo := model.GoogleUserInfo{
		ID:      "104829104857291039482",
		Email:   "dewi.lestari@gmail.com",
		Name:    "Dewi Lestari",
		Picture: "https://lh3.googleusercontent.com/a/ACg8ocP12Xz_dewi_avatar=s96-c",
	}

	mockOAuthClient.Mock.On("GetUserInfo", mock.Anything, "4/0AWgavTi-real_code_save_new_err_111").Return(userInfo, nil).Once()
	mockUserRepo.Mock.On("FindByEmail", mock.Anything, userInfo.Email).Return(entity.User{}, gorm.ErrRecordNotFound).Once()
	mockUserRepo.Mock.On("Save", mock.Anything, mock.Anything).Return(errors.New("Error 1062 (23000): Duplicate entry 'dewi.lestari@gmail.com' for key 'users.email'")).Once()

	request := httptest.NewRequest(fiber.MethodGet, "/api/auth/callback?code=4/0AWgavTi-real_code_save_new_err_111", nil)
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
}

func TestGoogleCallbackDatabaseErrorUpdateExistingUser(t *testing.T) {
	userInfo := model.GoogleUserInfo{
		ID:      "105938204918273645192",
		Email:   "eko.prasetyo@gmail.com",
		Name:    "Eko Prasetyo",
		Picture: "https://lh3.googleusercontent.com/a/ACg8ocR34Yw_eko_avatar=s96-c",
	}
	user := entity.User{
		ID:       "d4c1a6f2-9a5c-4c3b-8f2b-6a9e3c4d5e6f",
		Email:    userInfo.Email,
		Name:     userInfo.Name,
		GoogleID: userInfo.ID,
		Picture:  userInfo.Picture,
	}

	mockOAuthClient.Mock.On("GetUserInfo", mock.Anything, "4/0AWgavTj-real_code_update_err_222").Return(userInfo, nil).Once()
	mockUserRepo.Mock.On("FindByEmail", mock.Anything, userInfo.Email).Return(user, nil).Once()
	mockUserRepo.Mock.On("Save", mock.Anything, user).Return(errors.New("Error 1205 (HY000): Lock wait timeout exceeded; try restarting transaction")).Once()

	request := httptest.NewRequest(fiber.MethodGet, "/api/auth/callback?code=4/0AWgavTj-real_code_update_err_222", nil)
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
}

func TestGoogleCallbackTokenGenerationError(t *testing.T) {
	userInfo := model.GoogleUserInfo{
		ID:      "106827394819283746102",
		Email:   "rina.wijaya@gmail.com",
		Name:    "Rina Wijaya",
		Picture: "https://lh3.googleusercontent.com/a/ACg8ocT56Vu_rina_avatar=s96-c",
	}

	mockOAuthClient.Mock.On("GetUserInfo", mock.Anything, "4/0AWgavTk-real_code_token_gen_err_333").Return(userInfo, nil).Once()
	mockUserRepo.Mock.On("FindByEmail", mock.Anything, userInfo.Email).Return(entity.User{}, gorm.ErrRecordNotFound).Once()
	mockUserRepo.Mock.On("Save", mock.Anything, mock.Anything).Return(nil).Once()
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, mock.Anything).Return("", errors.New("key is of invalid type")).Once()

	request := httptest.NewRequest(fiber.MethodGet, "/api/auth/callback?code=4/0AWgavTk-real_code_token_gen_err_333", nil)
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
}

func TestGoogleCallbackDatabaseErrorSaveRefreshToken(t *testing.T) {
	userInfo := model.GoogleUserInfo{
		ID:      "107938475819283746591",
		Email:   "rizky.pratama@gmail.com",
		Name:    "Rizky Pratama",
		Picture: "https://lh3.googleusercontent.com/a/ACg8ocV78Ts_rizky_avatar=s96-c",
	}
	user := entity.User{
		ID:       "e5d2b7a3-0b6d-4d4c-9a3c-7b0f4d5e6f7a",
		Email:    userInfo.Email,
		Name:     userInfo.Name,
		GoogleID: userInfo.ID,
		Picture:  userInfo.Picture,
	}

	mockOAuthClient.Mock.On("GetUserInfo", mock.Anything, "4/0AWgavTl-real_code_token_save_err_444").Return(userInfo, nil).Once()
	mockUserRepo.Mock.On("FindByEmail", mock.Anything, userInfo.Email).Return(user, nil).Once()
	mockUserRepo.Mock.On("Save", mock.Anything, user).Return(nil).Once()
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, user).Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJlNWQyYjdhMy0wYjZkLTRkNGMtOWEzYy03YjBmNGQ1ZTZmN2EiLCJlbWFpbCI6InJpemt5LnByYXRhbWFAZ21haWwuY29tIiwiZXhwIjoxNzYwMTQ0MDAwfQ.sig-rizky-access-token", nil).Once()
	mockRefreshTokenRepo.Mock.On("Save", mock.Anything, mock.Anything).Return(errors.New("Error 1048 (23000): Column 'token' cannot be null")).Once()

	request := httptest.NewRequest(fiber.MethodGet, "/api/auth/callback?code=4/0AWgavTl-real_code_token_save_err_444", nil)
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
}


