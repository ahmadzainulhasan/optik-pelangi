package test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/config"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/entity"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

func TestGoogleLogin(t *testing.T) {
	mockOAuthClient.Mock.On("GetLoginURL", mock.Anything, mock.Anything).Return("https://accounts.google.com/o/oauth2/v2/auth?access_type=offline&client_id=718293847561-abcd1234efgh5678.apps.googleusercontent.com&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fapi%2Fauth%2Fcallback&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.profile&state=state-token")

	request := httptest.NewRequest(fiber.MethodGet, "/auth/login", nil)
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusSeeOther, response.StatusCode)
	assert.Equal(t, "https://accounts.google.com/o/oauth2/v2/auth?access_type=offline&client_id=718293847561-abcd1234efgh5678.apps.googleusercontent.com&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fapi%2Fauth%2Fcallback&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.profile&state=state-token", response.Header.Get("Location"))
	assert.Len(t, response.Cookies(), 2)
	assert.Equal(t, "oauth_state", response.Cookies()[0].Name)
	assert.Equal(t, "/auth/callback", response.Cookies()[0].Path)
	assert.Equal(t, http.SameSiteLaxMode, response.Cookies()[0].SameSite)
	assert.True(t, response.Cookies()[0].HttpOnly)
	assert.True(t, response.Cookies()[0].Secure)
	assert.Equal(t, "pkce_verifier", response.Cookies()[1].Name)
	assert.Equal(t, "/auth/callback", response.Cookies()[1].Path)
	assert.Equal(t, http.SameSiteLaxMode, response.Cookies()[1].SameSite)
	assert.True(t, response.Cookies()[1].HttpOnly)
	assert.True(t, response.Cookies()[1].Secure)
}

func TestGoogleCallbackSuccessNewUser(t *testing.T) {
	userInfo := model.GoogleUserInfo{
		ID:      "108234918237491823741",
		Email:   "budi.santoso@gmail.com",
		Name:    "Budi Santoso",
		Picture: "https://lh3.googleusercontent.com/a/ACg8ocIq7Q8m_budi_avatar=s96-c",
	}

	mockOAuthClient.Mock.On("GetUserInfo", mock.Anything, "4/0AWgavTe-real_code_new_user_xyz123", "valid-verifier-secret-1234567890").Return(userInfo, nil)
	mockUserRepo.Mock.On("FindByEmail", mock.Anything, userInfo.Email).Return(entity.User{}, gorm.ErrRecordNotFound)
	mockUserRepo.Mock.On("Save", mock.Anything, mock.MatchedBy(func(u entity.User) bool {
		return u.Email == "budi.santoso@gmail.com"
	})).Return(nil)
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, mock.MatchedBy(func(u entity.User) bool {
		return u.Email == "budi.santoso@gmail.com"
	})).Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJiNWYwZDNhMS1kZWIzLTQzNmQtOTYyMS02YzM1NWI3MjMwMWEiLCJlbWFpbCI6ImJ1ZGkuc2FudG9zb0BnbWFpbC5jb20iLCJleHAiOjE3NjAxNDQwMDB9.W-9vUfxlGk4_dK2M4NpXQyT8wZv5oP1aLmN3qRtUvY2", nil)
	mockRefreshTokenRepo.Mock.On("Save", mock.Anything, mock.MatchedBy(func(t entity.RefreshToken) bool {
		return t.UserID != "" && t.UserID != "c3b0f5e1-8f4b-4b2a-9e1a-5f8d2b3c4d5e" && t.UserID != "e5d2b7a3-0b6d-4d4c-9a3c-7b0f4d5e6f7a" && t.ID != "4c5d6e7f-8a9b-4c0d-1e2f-3a4b5c6d7e8f" && t.ID != "8b3f124c-9f8e-4a62-b91c-7d5a3e1b0c24" && t.ID != "5d6e7f8a-9b0c-4d1e-2f3a-4b5c6d7e8f9a"
	})).Return(nil)

	request := httptest.NewRequest(fiber.MethodGet, "/auth/callback?code=4/0AWgavTe-real_code_new_user_xyz123&state=valid-state-secret-1234567890", nil)
	request.Header.Set("Cookie", "oauth_state=valid-state-secret-1234567890; pkce_verifier=valid-verifier-secret-1234567890")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[model.LoginResponse]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJiNWYwZDNhMS1kZWIzLTQzNmQtOTYyMS02YzM1NWI3MjMwMWEiLCJlbWFpbCI6ImJ1ZGkuc2FudG9zb0BnbWFpbC5jb20iLCJleHAiOjE3NjAxNDQwMDB9.W-9vUfxlGk4_dK2M4NpXQyT8wZv5oP1aLmN3qRtUvY2", webResponse.Data.AccessToken)

	assert.Len(t, response.Cookies(), 3)

	var refreshTokenCookie, stateCookie, verifierCookie *http.Cookie
	for _, c := range response.Cookies() {
		switch c.Name {
		case "refresh_token":
			refreshTokenCookie = c
		case "oauth_state":
			stateCookie = c
		case "pkce_verifier":
			verifierCookie = c
		}
	}

	assert.NotNil(t, refreshTokenCookie)
	assert.Equal(t, "/auth", refreshTokenCookie.Path)
	assert.Equal(t, http.SameSiteStrictMode, refreshTokenCookie.SameSite)
	assert.True(t, refreshTokenCookie.Secure)
	assert.True(t, refreshTokenCookie.HttpOnly)

	assert.NotNil(t, stateCookie)
	assert.Equal(t, "", stateCookie.Value)
	assert.Equal(t, "/auth/callback", stateCookie.Path)

	assert.NotNil(t, verifierCookie)
	assert.Equal(t, "", verifierCookie.Value)
	assert.Equal(t, "/auth/callback", verifierCookie.Path)
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

	mockOAuthClient.Mock.On("GetUserInfo", mock.Anything, "4/0AWgavTf-real_code_existing_user_abc456", "valid-verifier-secret-1234567890").Return(userInfo, nil)
	mockUserRepo.Mock.On("FindByEmail", mock.Anything, userInfo.Email).Return(user, nil)
	mockUserRepo.Mock.On("Save", mock.Anything, user).Return(nil)
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, user).Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJjM2IwZjVlMS04ZjRiLTRiMmEtOWUxYS01ZjhkMmIzYzRkNWUiLCJlbWFpbCI6InNpdGkuYW1pbmFoQGdtYWlsLmNvbSIsImV4cCI6MTc2MDE0NDAwMH0.kK2M4NpXQyT8wZv5oP1aLmN3qRtUvY2jXk0W-9vUfxl", nil)
	mockRefreshTokenRepo.Mock.On("Save", mock.Anything, mock.MatchedBy(func(t entity.RefreshToken) bool {
		return t.UserID == "c3b0f5e1-8f4b-4b2a-9e1a-5f8d2b3c4d5e"
	})).Return(nil)

	request := httptest.NewRequest(fiber.MethodGet, "/auth/callback?code=4/0AWgavTf-real_code_existing_user_abc456&state=valid-state-secret-1234567890", nil)
	request.Header.Set("Cookie", "oauth_state=valid-state-secret-1234567890; pkce_verifier=valid-verifier-secret-1234567890")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[model.LoginResponse]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJjM2IwZjVlMS04ZjRiLTRiMmEtOWUxYS01ZjhkMmIzYzRkNWUiLCJlbWFpbCI6InNpdGkuYW1pbmFoQGdtYWlsLmNvbSIsImV4cCI6MTc2MDE0NDAwMH0.kK2M4NpXQyT8wZv5oP1aLmN3qRtUvY2jXk0W-9vUfxl", webResponse.Data.AccessToken)

	assert.Len(t, response.Cookies(), 3)

	var refreshTokenCookie, stateCookie, verifierCookie *http.Cookie
	for _, c := range response.Cookies() {
		switch c.Name {
		case "refresh_token":
			refreshTokenCookie = c
		case "oauth_state":
			stateCookie = c
		case "pkce_verifier":
			verifierCookie = c
		}
	}

	assert.NotNil(t, refreshTokenCookie)
	assert.Equal(t, "/auth", refreshTokenCookie.Path)
	assert.Equal(t, http.SameSiteStrictMode, refreshTokenCookie.SameSite)
	assert.True(t, refreshTokenCookie.Secure)
	assert.True(t, refreshTokenCookie.HttpOnly)

	assert.NotNil(t, stateCookie)
	assert.Equal(t, "", stateCookie.Value)
	assert.Equal(t, "/auth/callback", stateCookie.Path)

	assert.NotNil(t, verifierCookie)
	assert.Equal(t, "", verifierCookie.Value)
	assert.Equal(t, "/auth/callback", verifierCookie.Path)
}

func TestGoogleCallbackStateMismatch(t *testing.T) {
	request := httptest.NewRequest(fiber.MethodGet, "/auth/callback?code=4/0AWgavT-attacker-code&state=forged-state-attack", nil)
	request.Header.Set("Cookie", "oauth_state=original-legitimate-state; pkce_verifier=legitimate-verifier")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[map[string]string]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, webResponse.Data)
	assert.Equal(t, "State harus sama dengan CookieState", webResponse.Data["State"])

	assert.Len(t, response.Cookies(), 2)
	assert.Equal(t, "", response.Cookies()[0].Value)
	assert.Equal(t, "", response.Cookies()[1].Value)
}

func TestGoogleCallbackMissingVerifier(t *testing.T) {
	request := httptest.NewRequest(fiber.MethodGet, "/auth/callback?code=4/0AWgavT-valid-code&state=valid-state", nil)
	request.Header.Set("Cookie", "oauth_state=valid-state")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[map[string]string]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, webResponse.Data)
	assert.Equal(t, "Verifier wajib diisi", webResponse.Data["Verifier"])
}

func TestGoogleCallbackMissingStateCookie(t *testing.T) {
	request := httptest.NewRequest(fiber.MethodGet, "/auth/callback?code=4/0AWgavT-valid-code&state=valid-state-123", nil)
	request.Header.Set("Cookie", "pkce_verifier=valid-verifier-123")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[map[string]string]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, webResponse.Data)
	assert.Equal(t, "CookieState wajib diisi", webResponse.Data["CookieState"])
}

func TestGoogleCallbackMissingQueryState(t *testing.T) {
	request := httptest.NewRequest(fiber.MethodGet, "/auth/callback?code=4/0AWgavT-valid-code", nil)
	request.Header.Set("Cookie", "oauth_state=valid-state-123; pkce_verifier=valid-verifier-123")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[map[string]string]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, webResponse.Data)
	assert.Equal(t, "State wajib diisi", webResponse.Data["State"])
}

func TestGoogleCallbackGoogleError(t *testing.T) {
	request := httptest.NewRequest(fiber.MethodGet, "/auth/callback?error=access_denied", nil)
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[map[string]string]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, webResponse.Data)
	assert.Equal(t, "Error harus kosong", webResponse.Data["Error"])
}

func TestGoogleCallbackMissingCode(t *testing.T) {
	request := httptest.NewRequest(fiber.MethodGet, "/auth/callback?state=valid-state-123", nil)
	request.Header.Set("Cookie", "oauth_state=valid-state-123; pkce_verifier=valid-verifier-123")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[map[string]string]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, webResponse.Data)
	assert.Equal(t, "Code wajib diisi", webResponse.Data["Code"])
}

func TestGoogleCallbackEmptyCode(t *testing.T) {
	request := httptest.NewRequest(fiber.MethodGet, "/auth/callback?code=&state=valid-state-123", nil)
	request.Header.Set("Cookie", "oauth_state=valid-state-123; pkce_verifier=valid-verifier-123")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[map[string]string]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, webResponse.Data)
	assert.Equal(t, "Code wajib diisi", webResponse.Data["Code"])
}

func TestGoogleCallbackOAuthExchangeFailed(t *testing.T) {
	mockOAuthClient.Mock.On("GetUserInfo", mock.Anything, "4/0AWgavTg-invalid-or-expired-auth-token-999", "valid-verifier-secret-1234567890").
		Return(model.GoogleUserInfo{}, &oauth2.RetrieveError{
			Response: &http.Response{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
			},
			ErrorCode:        "invalid_grant",
			ErrorDescription: "Bad Request",
		})

	request := httptest.NewRequest(fiber.MethodGet, "/auth/callback?code=4/0AWgavTg-invalid-or-expired-auth-token-999&state=valid-state-secret-1234567890", nil)
	request.Header.Set("Cookie", "oauth_state=valid-state-secret-1234567890; pkce_verifier=valid-verifier-secret-1234567890")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[any]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, webResponse.Errors)
}

func TestGoogleCallbackDatabaseErrorFindUser(t *testing.T) {
	userInfo := model.GoogleUserInfo{
		ID:      "109823749182374918234",
		Email:   "ahmad.dahlan@gmail.com",
		Name:    "Ahmad Dahlan",
		Picture: "https://lh3.googleusercontent.com/a/ACg8ocM93Lk_ahmad_avatar=s96-c",
	}

	mockOAuthClient.Mock.On("GetUserInfo", mock.Anything, "4/0AWgavTh-real_code_find_err_789", "valid-verifier-secret-1234567890").Return(userInfo, nil)
	mockUserRepo.Mock.On("FindByEmail", mock.Anything, userInfo.Email).Return(entity.User{}, gorm.ErrInvalidDB)

	request := httptest.NewRequest(fiber.MethodGet, "/auth/callback?code=4/0AWgavTh-real_code_find_err_789&state=valid-state-secret-1234567890", nil)
	request.Header.Set("Cookie", "oauth_state=valid-state-secret-1234567890; pkce_verifier=valid-verifier-secret-1234567890")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[any]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, webResponse.Errors)
}

func TestGoogleCallbackDatabaseErrorSaveNewUser(t *testing.T) {
	userInfo := model.GoogleUserInfo{
		ID:      "104829104857291039482",
		Email:   "dewi.lestari@gmail.com",
		Name:    "Dewi Lestari",
		Picture: "https://lh3.googleusercontent.com/a/ACg8ocP12Xz_dewi_avatar=s96-c",
	}

	mockOAuthClient.Mock.On("GetUserInfo", mock.Anything, "4/0AWgavTi-real_code_save_new_err_111", "valid-verifier-secret-1234567890").Return(userInfo, nil)
	mockUserRepo.Mock.On("FindByEmail", mock.Anything, userInfo.Email).Return(entity.User{}, gorm.ErrRecordNotFound)
	mockUserRepo.Mock.On("Save", mock.Anything, mock.MatchedBy(func(u entity.User) bool {
		return u.Email == "dewi.lestari@gmail.com"
	})).Return(gorm.ErrDuplicatedKey)

	request := httptest.NewRequest(fiber.MethodGet, "/auth/callback?code=4/0AWgavTi-real_code_save_new_err_111&state=valid-state-secret-1234567890", nil)
	request.Header.Set("Cookie", "oauth_state=valid-state-secret-1234567890; pkce_verifier=valid-verifier-secret-1234567890")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[any]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, webResponse.Errors)
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

	mockOAuthClient.Mock.On("GetUserInfo", mock.Anything, "4/0AWgavTj-real_code_update_err_222", "valid-verifier-secret-1234567890").Return(userInfo, nil)
	mockUserRepo.Mock.On("FindByEmail", mock.Anything, userInfo.Email).Return(user, nil)
	mockUserRepo.Mock.On("Save", mock.Anything, user).Return(gorm.ErrInvalidTransaction)

	request := httptest.NewRequest(fiber.MethodGet, "/auth/callback?code=4/0AWgavTj-real_code_update_err_222&state=valid-state-secret-1234567890", nil)
	request.Header.Set("Cookie", "oauth_state=valid-state-secret-1234567890; pkce_verifier=valid-verifier-secret-1234567890")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[any]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, webResponse.Errors)
}

func TestGoogleCallbackTokenGenerationError(t *testing.T) {
	userInfo := model.GoogleUserInfo{
		ID:      "106827394819283746102",
		Email:   "rina.wijaya@gmail.com",
		Name:    "Rina Wijaya",
		Picture: "https://lh3.googleusercontent.com/a/ACg8ocT56Vu_rina_avatar=s96-c",
	}

	user := entity.User{
		ID:       "d5e6f7a8-b9c0-4d1e-2f3a-4b5c6d7e8f9a",
		Email:    userInfo.Email,
		Name:     userInfo.Name,
		GoogleID: userInfo.ID,
		Picture:  userInfo.Picture,
	}

	mockOAuthClient.Mock.On("GetUserInfo", mock.Anything, "4/0AWgavTk-real_code_token_gen_err_333", "valid-verifier-secret-1234567890").Return(userInfo, nil)
	mockUserRepo.Mock.On("FindByEmail", mock.Anything, userInfo.Email).Return(user, nil)
	mockUserRepo.Mock.On("Save", mock.Anything, user).Return(nil)
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, user).Return("", jwt.ErrInvalidKeyType)

	request := httptest.NewRequest(fiber.MethodGet, "/auth/callback?code=4/0AWgavTk-real_code_token_gen_err_333&state=valid-state-secret-1234567890", nil)
	request.Header.Set("Cookie", "oauth_state=valid-state-secret-1234567890; pkce_verifier=valid-verifier-secret-1234567890")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[any]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, webResponse.Errors)
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

	mockOAuthClient.Mock.On("GetUserInfo", mock.Anything, "4/0AWgavTl-real_code_token_save_err_444", "valid-verifier-secret-1234567890").Return(userInfo, nil)
	mockUserRepo.Mock.On("FindByEmail", mock.Anything, userInfo.Email).Return(user, nil)
	mockUserRepo.Mock.On("Save", mock.Anything, user).Return(nil)
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, user).Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJlNWQyYjdhMy0wYjZkLTRkNGMtOWEzYy03YjBmNGQ1ZTZmN2EiLCJlbWFpbCI6InJpemt5LnByYXRhbWFAZ21haWwuY29tIiwiZXhwIjoxNzYwMTQ0MDAwfQ.v5oP1aLmN3qRtUvY2jXk0W-9vUfxlGk4_dK2M4NpXQy", nil)
	mockRefreshTokenRepo.Mock.On("Save", mock.Anything, mock.MatchedBy(func(t entity.RefreshToken) bool {
		return t.UserID == "e5d2b7a3-0b6d-4d4c-9a3c-7b0f4d5e6f7a"
	})).Return(gorm.ErrInvalidData)

	request := httptest.NewRequest(fiber.MethodGet, "/auth/callback?code=4/0AWgavTl-real_code_token_save_err_444&state=valid-state-secret-1234567890", nil)
	request.Header.Set("Cookie", "oauth_state=valid-state-secret-1234567890; pkce_verifier=valid-verifier-secret-1234567890")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[any]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, webResponse.Errors)
}


func TestRefreshTokenSuccess(t *testing.T) {
	user := entity.User{
		ID:      "a1b2c3d4-e5f6-4a1b-8c2d-3e4f5a6b7c8d",
		Email:   "siti.rahmawati@gmail.com",
		Name:    "Siti Rahmawati",
		Picture: "https://lh3.googleusercontent.com/a/ACg8ocIq7Q8m_siti_avatar=s96-c",
	}

	refreshToken := entity.RefreshToken{
		ID:        "8b3f124c-9f8e-4a62-b91c-7d5a3e1b0c24",
		UserID:    user.ID,
		Token:     "c4b8e21a-7f39-4d65-8b12-9e8c3a1f4b5d",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	mockRefreshTokenRepo.Mock.On("FindByToken", mock.Anything, "c4b8e21a-7f39-4d65-8b12-9e8c3a1f4b5d").Return(refreshToken, nil)
	mockUserRepo.Mock.On("FindByID", mock.Anything, user.ID).Return(user, nil)
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, user).Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhMWIyYzNkNC1lNWY2LTRhMWItOGMyZC0zZTRmNWE2YjdjOGQiLCJlbWFpbCI6InNpdGkucmFobWF3YXRpQGdtYWlsLmNvbSIsImV4cCI6MTc2MDE0NDAwMH0.qRtUvY2jXk1i7Qx9pZkK2M4NpXQyT8wZvwN5yP8xLm0", nil)
	mockRefreshTokenRepo.Mock.On("Save", mock.Anything, mock.MatchedBy(func(t entity.RefreshToken) bool {
		return t.UserID == "a1b2c3d4-e5f6-4a1b-8c2d-3e4f5a6b7c8d"
	})).Return(nil)

	request := httptest.NewRequest(fiber.MethodPost, "/auth/refresh", nil)
	request.Header.Set("Cookie", "refresh_token=c4b8e21a-7f39-4d65-8b12-9e8c3a1f4b5d")
	request.Header.Set("Authorization", "Bearer old-access-token-to-revoke-12345")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[model.LoginResponse]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhMWIyYzNkNC1lNWY2LTRhMWItOGMyZC0zZTRmNWE2YjdjOGQiLCJlbWFpbCI6InNpdGkucmFobWF3YXRpQGdtYWlsLmNvbSIsImV4cCI6MTc2MDE0NDAwMH0.qRtUvY2jXk1i7Qx9pZkK2M4NpXQyT8wZvwN5yP8xLm0", webResponse.Data.AccessToken)

	assert.Len(t, response.Cookies(), 1)
	assert.Equal(t, "refresh_token", response.Cookies()[0].Name)
	assert.Equal(t, "/auth", response.Cookies()[0].Path)
	assert.Equal(t, http.SameSiteStrictMode, response.Cookies()[0].SameSite)
	assert.True(t, response.Cookies()[0].HttpOnly)
	assert.True(t, response.Cookies()[0].Secure)
	assert.NotEmpty(t, response.Cookies()[0].Value)
	assert.NotEqual(t, "c4b8e21a-7f39-4d65-8b12-9e8c3a1f4b5d", response.Cookies()[0].Value)
}

func TestRefreshTokenUnauthorizedNoCookie(t *testing.T) {
	request := httptest.NewRequest(fiber.MethodPost, "/auth/refresh", nil)
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
}

func TestRefreshTokenUnauthorizedEmptyCookie(t *testing.T) {
	request := httptest.NewRequest(fiber.MethodPost, "/auth/refresh", nil)
	request.Header.Set("Cookie", "refresh_token=   ")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
}

func TestRefreshTokenInvalidOrExpired(t *testing.T) {
	mockRefreshTokenRepo.Mock.On("FindByToken", mock.Anything, "e2a7b8c9-3d1f-4a5e-b6c8-9d0e1f2a3b4c").Return(entity.RefreshToken{}, gorm.ErrRecordNotFound)

	request := httptest.NewRequest(fiber.MethodPost, "/auth/refresh", nil)
	request.Header.Set("Cookie", "refresh_token=e2a7b8c9-3d1f-4a5e-b6c8-9d0e1f2a3b4c")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
	assert.NotEmpty(t, response.Cookies())
}

func TestRefreshTokenDatabaseErrorFindByToken(t *testing.T) {
	mockRefreshTokenRepo.Mock.On("FindByToken", mock.Anything, "7f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b5c").Return(entity.RefreshToken{}, errors.New("db find failed"))

	request := httptest.NewRequest(fiber.MethodPost, "/auth/refresh", nil)
	request.Header.Set("Cookie", "refresh_token=7f1a2b3c-4d5e-6f7a-8b9c-0d1e2f3a4b5c")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
	assert.Empty(t, response.Cookies())
}

func TestRefreshTokenDatabaseErrorFindUser(t *testing.T) {
	refreshToken := entity.RefreshToken{
		ID:        "1f2a3b4c-5d6e-7f8a-9b0c-1d2e3f4a5b6c",
		UserID:    "9a8b7c6d-5e4f-3a2b-1c0d-e9f8a7b6c5d4",
		Token:     "a3b4c5d6-e7f8-4a9b-8c1d-2e3f4a5b6c7d",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	mockRefreshTokenRepo.Mock.On("FindByToken", mock.Anything, "a3b4c5d6-e7f8-4a9b-8c1d-2e3f4a5b6c7d").Return(refreshToken, nil)
	mockUserRepo.Mock.On("FindByID", mock.Anything, "9a8b7c6d-5e4f-3a2b-1c0d-e9f8a7b6c5d4").Return(entity.User{}, errors.New("db find user error"))

	request := httptest.NewRequest(fiber.MethodPost, "/auth/refresh", nil)
	request.Header.Set("Cookie", "refresh_token=a3b4c5d6-e7f8-4a9b-8c1d-2e3f4a5b6c7d")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
}

func TestRefreshTokenExpiredToken(t *testing.T) {
	refreshToken := entity.RefreshToken{
		ID:        "2a3b4c5d-6e7f-4a8b-9c0d-1e2f3a4b5c6d",
		UserID:    "b4c5d6e7-f8a9-4b0c-9d2e-3f4a5b6c7d8e",
		Token:     "d6e7f8a9-b0c1-4d2e-8f3a-4b5c6d7e8f9a",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	mockRefreshTokenRepo.Mock.On("FindByToken", mock.Anything, "d6e7f8a9-b0c1-4d2e-8f3a-4b5c6d7e8f9a").Return(refreshToken, nil)

	request := httptest.NewRequest(fiber.MethodPost, "/auth/refresh", nil)
	request.Header.Set("Cookie", "refresh_token=d6e7f8a9-b0c1-4d2e-8f3a-4b5c6d7e8f9a")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
}

func TestRefreshTokenTokenGenerationError(t *testing.T) {
	user := entity.User{
		ID:    "c3d4e5f6-a1b2-4c3d-ae4f-5a6b7c8d9e0f",
		Email: "putri.anggraini@gmail.com",
	}
	refreshToken := entity.RefreshToken{
		ID:        "3b4c5d6e-7f8a-4b9c-0d1e-2f3a4b5c6d7e",
		UserID:    user.ID,
		Token:     "e7f8a9b0-c1d2-4e3f-9a0b-1c2d3e4f5a6b",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	mockRefreshTokenRepo.Mock.On("FindByToken", mock.Anything, "e7f8a9b0-c1d2-4e3f-9a0b-1c2d3e4f5a6b").Return(refreshToken, nil)
	mockUserRepo.Mock.On("FindByID", mock.Anything, user.ID).Return(user, nil)
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, user).Return("", errors.New("jwt sign error"))

	request := httptest.NewRequest(fiber.MethodPost, "/auth/refresh", nil)
	request.Header.Set("Cookie", "refresh_token=e7f8a9b0-c1d2-4e3f-9a0b-1c2d3e4f5a6b")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
}

func TestRefreshTokenDatabaseErrorSaveNewToken(t *testing.T) {
	user := entity.User{
		ID:    "f6a1b2c3-4d4e-bf5a-6b7c-8d9e0f1a2b3c",
		Email: "bayu.wicaksono@gmail.com",
	}
	refreshToken := entity.RefreshToken{
		ID:        "4c5d6e7f-8a9b-4c0d-1e2f-3a4b5c6d7e8f",
		UserID:    user.ID,
		Token:     "f8a9b0c1-d2e3-4f0a-8b1c-2d3e4f5a6b7c",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	mockRefreshTokenRepo.Mock.On("FindByToken", mock.Anything, "f8a9b0c1-d2e3-4f0a-8b1c-2d3e4f5a6b7c").Return(refreshToken, nil)
	mockUserRepo.Mock.On("FindByID", mock.Anything, user.ID).Return(user, nil)
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, user).Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmNmExYjJjMy00ZDRlLWJmNWEtNmI3Yy04ZDllMGYxYTJiM2MiLCJlbWFpbCI6ImJheXUud2ljYWtzb25vQGdtYWlsLmNvbSIsImV4cCI6MTc2MDE0NDAwMH0.i7Qx9pZkK2M4NpXQyT8wZvwN5yP8xLm0qRtUvY2jXk1", nil)
	mockRefreshTokenRepo.Mock.On("Save", mock.Anything, mock.MatchedBy(func(t entity.RefreshToken) bool {
		return t.ID == "4c5d6e7f-8a9b-4c0d-1e2f-3a4b5c6d7e8f"
	})).Return(errors.New("db save failed"))

	request := httptest.NewRequest(fiber.MethodPost, "/auth/refresh", nil)
	request.Header.Set("Cookie", "refresh_token=f8a9b0c1-d2e3-4f0a-8b1c-2d3e4f5a6b7c")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
}

func TestRefreshTokenGracePeriodFromRedis(t *testing.T) {
	redisClient := config.NewRedis(viperConfig, log)
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		t.Skip("redis is not available")
	}

	testUsecase := usecase.NewAuthUsecaseImpl(
		db,
		redisClient,
		mockUserRepo,
		mockRefreshTokenRepo,
		mockJwtUtil,
		validate,
		log,
		mockOAuthClient,
	)

	user := entity.User{
		ID:    "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
		Email: "grace.period@gmail.com",
	}

	token1 := uuid.Must(uuid.NewV7()).String()
	refreshToken := entity.RefreshToken{
		ID:        "5d6e7f8a-9b0c-4d1e-2f3a-4b5c6d7e8f9a",
		UserID:    user.ID,
		Token:     token1,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	mockRefreshTokenRepo.Mock.On("FindByToken", mock.Anything, token1).Return(refreshToken, nil)
	mockUserRepo.Mock.On("FindByID", mock.Anything, user.ID).Return(user, nil)
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, user).Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhMWIyYzNkNC1lNWY2LTc4OTAtYWJjZC1lZjEyMzQ1Njc4OTAiLCJlbWFpbCI6ImdyYWNlLnBlcmlvZEBnbWFpbC5jb20iLCJleHAiOjE3NjAxNDQwMDB9.wZvwN5yP8xLm0qRtUvY2jXk1i7Qx9pZkK2M4NpXQyT8", nil)
	mockRefreshTokenRepo.Mock.On("Save", mock.Anything, mock.MatchedBy(func(t entity.RefreshToken) bool {
		return t.UserID == "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	})).Return(nil)

	res1, err := testUsecase.RotateToken(context.Background(), token1)
	assert.NoError(t, err)
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhMWIyYzNkNC1lNWY2LTc4OTAtYWJjZC1lZjEyMzQ1Njc4OTAiLCJlbWFpbCI6ImdyYWNlLnBlcmlvZEBnbWFpbC5jb20iLCJleHAiOjE3NjAxNDQwMDB9.wZvwN5yP8xLm0qRtUvY2jXk1i7Qx9pZkK2M4NpXQyT8", res1.AccessToken)
	assert.NotEmpty(t, res1.RefreshToken)

	res2, err := testUsecase.RotateToken(context.Background(), token1)
	assert.NoError(t, err)
	assert.Equal(t, res1.AccessToken, res2.AccessToken)
	assert.Equal(t, res1.RefreshToken, res2.RefreshToken)
}

func TestLogoutSuccess(t *testing.T) {
	refreshToken := "e4d909c2-9012-49e8-b769-3c9f2d123456"
	mockRefreshTokenRepo.Mock.On("DeleteByToken", mock.Anything, refreshToken).Return(nil)

	request := httptest.NewRequest(fiber.MethodDelete, "/auth/logout", nil)
	request.Header.Set("Cookie", "refresh_token="+refreshToken)
	request.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJjM2IwZjVlMS04ZjRiLTRiMmEtOWUxYS01ZjhkMmIzYzRkNWUiLCJleHAiOjE3NjAxNDQwMDB9.kK2M4NpXQyT8wZv5oP1aLmN3qRtUvY2jXk0W-9vUfxl")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[string]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.Equal(t, "OK", webResponse.Data)

	var refreshTokenCookie *http.Cookie
	for _, c := range response.Cookies() {
		if c.Name == "refresh_token" {
			refreshTokenCookie = c
			break
		}
	}

	assert.NotNil(t, refreshTokenCookie)
	assert.Equal(t, "", refreshTokenCookie.Value)
	assert.Equal(t, "/auth", refreshTokenCookie.Path)
	assert.True(t, refreshTokenCookie.HttpOnly)
	assert.True(t, refreshTokenCookie.Secure)
	assert.Equal(t, http.SameSiteStrictMode, refreshTokenCookie.SameSite)
}

func TestLogoutSuccessNoTokens(t *testing.T) {
	mockRefreshTokenRepo.Mock.On("DeleteByToken", mock.Anything, "").Return(nil)

	request := httptest.NewRequest(fiber.MethodDelete, "/auth/logout", nil)
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[string]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.Equal(t, "OK", webResponse.Data)

	var refreshTokenCookie *http.Cookie
	for _, c := range response.Cookies() {
		if c.Name == "refresh_token" {
			refreshTokenCookie = c
			break
		}
	}

	assert.NotNil(t, refreshTokenCookie)
	assert.Equal(t, "", refreshTokenCookie.Value)
	assert.Equal(t, "/auth", refreshTokenCookie.Path)
}

func TestLogoutSuccessOnlyBearerToken(t *testing.T) {
	mockRefreshTokenRepo.Mock.On("DeleteByToken", mock.Anything, "").Return(nil)

	request := httptest.NewRequest(fiber.MethodDelete, "/auth/logout", nil)
	request.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJiNWYwZDNhMS1kZWIzLTQzNmQtOTYyMS02YzM1NWI3MjMwMWEiLCJleHAiOjE3NjAxNDQwMDB9.W-9vUfxlGk4_dK2M4NpXQyT8wZv5oP1aLmN3qRtUvY2")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[string]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.Equal(t, "OK", webResponse.Data)

	var refreshTokenCookie *http.Cookie
	for _, c := range response.Cookies() {
		if c.Name == "refresh_token" {
			refreshTokenCookie = c
			break
		}
	}

	assert.NotNil(t, refreshTokenCookie)
	assert.Equal(t, "", refreshTokenCookie.Value)
	assert.Equal(t, "/auth", refreshTokenCookie.Path)
}

func TestLogoutDatabaseError(t *testing.T) {
	refreshToken := "f81d4fae-7dec-11d0-a765-00a0c91e6bf6"
	mockRefreshTokenRepo.Mock.On("DeleteByToken", mock.Anything, refreshToken).Return(errors.New("connection reset by peer"))

	request := httptest.NewRequest(fiber.MethodDelete, "/auth/logout", nil)
	request.Header.Set("Cookie", "refresh_token="+refreshToken)
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[string]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.Equal(t, "OK", webResponse.Data)

	var refreshTokenCookie *http.Cookie
	for _, c := range response.Cookies() {
		if c.Name == "refresh_token" {
			refreshTokenCookie = c
			break
		}
	}

	assert.NotNil(t, refreshTokenCookie)
	assert.Equal(t, "", refreshTokenCookie.Value)
	assert.Equal(t, "/auth", refreshTokenCookie.Path)
}
