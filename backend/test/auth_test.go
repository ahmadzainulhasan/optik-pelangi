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
	assert.Equal(t, "/", response.Cookies()[0].Path)
	assert.Equal(t, http.SameSiteLaxMode, response.Cookies()[0].SameSite)
	assert.True(t, response.Cookies()[0].HttpOnly)
	assert.True(t, response.Cookies()[0].Secure)
	assert.Equal(t, "pkce_verifier", response.Cookies()[1].Name)
	assert.Equal(t, "/", response.Cookies()[1].Path)
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
		return u.Email == userInfo.Email
	})).Return(nil)
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, mock.MatchedBy(func(u entity.User) bool {
		return u.Email == userInfo.Email
	})).Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJiNWYwZDNhMS1kZWIzLTQzNmQtOTYyMS02YzM1NWI3MjMwMWEiLCJlbWFpbCI6ImJ1ZGkuc2FudG9zb0BnbWFpbC5jb20iLCJleHAiOjE3NjAxNDQwMDB9.sig-budi-access-token", nil)
	mockRefreshTokenRepo.Mock.On("Save", mock.Anything, mock.MatchedBy(func(rt entity.RefreshToken) bool {
		return rt.UserID != "c3b0f5e1-8f4b-4b2a-9e1a-5f8d2b3c4d5e" && rt.UserID != "e5d2b7a3-0b6d-4d4c-9a3c-7b0f4d5e6f7a" && rt.UserID != "f6a1b2c3-4d4e-bf5a-6b7c-8d9e0f1a2b3c"
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
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJiNWYwZDNhMS1kZWIzLTQzNmQtOTYyMS02YzM1NWI3MjMwMWEiLCJlbWFpbCI6ImJ1ZGkuc2FudG9zb0BnbWFpbC5jb20iLCJleHAiOjE3NjAxNDQwMDB9.sig-budi-access-token", webResponse.Data.AccessToken)

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

	assert.NotNil(t, verifierCookie)
	assert.Equal(t, "", verifierCookie.Value)
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
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, mock.MatchedBy(func(u entity.User) bool {
		return u.Email == userInfo.Email
	})).Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJjM2IwZjVlMS04ZjRiLTRiMmEtOWUxYS01ZjhkMmIzYzRkNWUiLCJlbWFpbCI6InNpdGkuYW1pbmFoQGdtYWlsLmNvbSIsImV4cCI6MTc2MDE0NDAwMH0.sig-siti-access-token", nil)
	mockRefreshTokenRepo.Mock.On("Save", mock.Anything, mock.MatchedBy(func(rt entity.RefreshToken) bool {
		return rt.UserID == user.ID
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
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJjM2IwZjVlMS04ZjRiLTRiMmEtOWUxYS01ZjhkMmIzYzRkNWUiLCJlbWFpbCI6InNpdGkuYW1pbmFoQGdtYWlsLmNvbSIsImV4cCI6MTc2MDE0NDAwMH0.sig-siti-access-token", webResponse.Data.AccessToken)

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

	assert.NotNil(t, verifierCookie)
	assert.Equal(t, "", verifierCookie.Value)
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
		return u.Email == userInfo.Email
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

	mockOAuthClient.Mock.On("GetUserInfo", mock.Anything, "4/0AWgavTk-real_code_token_gen_err_333", "valid-verifier-secret-1234567890").Return(userInfo, nil)
	mockUserRepo.Mock.On("FindByEmail", mock.Anything, userInfo.Email).Return(entity.User{}, gorm.ErrRecordNotFound)
	mockUserRepo.Mock.On("Save", mock.Anything, mock.MatchedBy(func(u entity.User) bool {
		return u.Email == userInfo.Email
	})).Return(nil)
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, mock.MatchedBy(func(u entity.User) bool {
		return u.Email == userInfo.Email
	})).Return("", jwt.ErrInvalidKeyType)

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
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, mock.MatchedBy(func(u entity.User) bool {
		return u.Email == userInfo.Email
	})).Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJlNWQyYjdhMy0wYjZkLTRkNGMtOWEzYy03YjBmNGQ1ZTZmN2EiLCJlbWFpbCI6InJpemt5LnByYXRhbWFAZ21haWwuY29tIiwiZXhwIjoxNzYwMTQ0MDAwfQ.sig-rizky-access-token", nil)
	mockRefreshTokenRepo.Mock.On("Save", mock.Anything, mock.MatchedBy(func(rt entity.RefreshToken) bool {
		return rt.UserID == user.ID
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

	mockRefreshTokenRepo.Mock.On("FindByToken", mock.Anything, "valid-refresh-token-in-db-111").Return(entity.RefreshToken{
		ID:        "token-id-1",
		UserID:    user.ID,
		Token:     "valid-refresh-token-in-db-111",
		ExpiresAt: time.Now().Add(time.Hour),
	}, nil)
	mockUserRepo.Mock.On("FindByID", mock.Anything, user.ID).Return(user, nil)
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, mock.MatchedBy(func(u entity.User) bool {
		return u.ID == user.ID
	})).Return("new-rotated-access-token-jwt", nil)
	mockRefreshTokenRepo.Mock.On("Save", mock.Anything, mock.MatchedBy(func(rt entity.RefreshToken) bool {
		return rt.ID == "token-id-1" && rt.UserID == user.ID && rt.Token != "valid-refresh-token-in-db-111"
	})).Return(nil)

	request := httptest.NewRequest(fiber.MethodPost, "/auth/refresh", nil)
	request.Header.Set("Cookie", "refresh_token=valid-refresh-token-in-db-111")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	var webResponse model.WebResponse[model.LoginResponse]
	err = json.Unmarshal(bodyBytes, &webResponse)
	assert.NoError(t, err)
	assert.Equal(t, "new-rotated-access-token-jwt", webResponse.Data.AccessToken)

	assert.Len(t, response.Cookies(), 1)
	assert.Equal(t, "refresh_token", response.Cookies()[0].Name)
	assert.Equal(t, "/auth", response.Cookies()[0].Path)
	assert.Equal(t, http.SameSiteStrictMode, response.Cookies()[0].SameSite)
	assert.True(t, response.Cookies()[0].HttpOnly)
	assert.True(t, response.Cookies()[0].Secure)
	assert.NotEmpty(t, response.Cookies()[0].Value)
	assert.NotEqual(t, "valid-refresh-token-in-db-111", response.Cookies()[0].Value)
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
	mockRefreshTokenRepo.Mock.On("FindByToken", mock.Anything, "expired-or-revoked-token-333").Return(entity.RefreshToken{}, gorm.ErrRecordNotFound)

	request := httptest.NewRequest(fiber.MethodPost, "/auth/refresh", nil)
	request.Header.Set("Cookie", "refresh_token=expired-or-revoked-token-333")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
	assert.NotEmpty(t, response.Cookies())
}

func TestRefreshTokenDatabaseErrorFindByToken(t *testing.T) {
	mockRefreshTokenRepo.Mock.On("FindByToken", mock.Anything, "db-err-find-token-222").Return(entity.RefreshToken{}, errors.New("db find failed"))

	request := httptest.NewRequest(fiber.MethodPost, "/auth/refresh", nil)
	request.Header.Set("Cookie", "refresh_token=db-err-find-token-222")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
	assert.Empty(t, response.Cookies())
}

func TestRefreshTokenDatabaseErrorFindUser(t *testing.T) {
	mockRefreshTokenRepo.Mock.On("FindByToken", mock.Anything, "valid-token-user-deleted-444").Return(entity.RefreshToken{
		ID:        "token-id-4",
		UserID:    "non-existent-user-uuid",
		Token:     "valid-token-user-deleted-444",
		ExpiresAt: time.Now().Add(time.Hour),
	}, nil)
	mockUserRepo.Mock.On("FindByID", mock.Anything, "non-existent-user-uuid").Return(entity.User{}, errors.New("db find user error"))

	request := httptest.NewRequest(fiber.MethodPost, "/auth/refresh", nil)
	request.Header.Set("Cookie", "refresh_token=valid-token-user-deleted-444")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
}

func TestRefreshTokenExpiredToken(t *testing.T) {
	mockRefreshTokenRepo.Mock.On("FindByToken", mock.Anything, "token-already-expired-555").Return(entity.RefreshToken{
		ID:        "token-id-5",
		UserID:    "user-id-5",
		Token:     "token-already-expired-555",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}, nil)

	request := httptest.NewRequest(fiber.MethodPost, "/auth/refresh", nil)
	request.Header.Set("Cookie", "refresh_token=token-already-expired-555")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
}

func TestRefreshTokenTokenGenerationError(t *testing.T) {
	user := entity.User{
		ID:    "c3d4e5f6-a1b2-4c3d-ae4f-5a6b7c8d9e0f",
		Email: "putri.anggraini@gmail.com",
	}

	mockRefreshTokenRepo.Mock.On("FindByToken", mock.Anything, "jwt-err-token-666").Return(entity.RefreshToken{
		ID:        "token-id-6",
		UserID:    user.ID,
		Token:     "jwt-err-token-666",
		ExpiresAt: time.Now().Add(time.Hour),
	}, nil)
	mockUserRepo.Mock.On("FindByID", mock.Anything, user.ID).Return(user, nil)
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, mock.MatchedBy(func(u entity.User) bool {
		return u.ID == user.ID
	})).Return("", errors.New("jwt sign error"))

	request := httptest.NewRequest(fiber.MethodPost, "/auth/refresh", nil)
	request.Header.Set("Cookie", "refresh_token=jwt-err-token-666")
	response, err := app.Test(request)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
}

func TestRefreshTokenDatabaseErrorSaveNewToken(t *testing.T) {
	user := entity.User{
		ID:    "f6a1b2c3-4d4e-bf5a-6b7c-8d9e0f1a2b3c",
		Email: "bayu.wicaksono@gmail.com",
	}

	mockRefreshTokenRepo.Mock.On("FindByToken", mock.Anything, "db-err-save-token-777").Return(entity.RefreshToken{
		ID:        "token-id-7",
		UserID:    user.ID,
		Token:     "db-err-save-token-777",
		ExpiresAt: time.Now().Add(time.Hour),
	}, nil)
	mockUserRepo.Mock.On("FindByID", mock.Anything, user.ID).Return(user, nil)
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, mock.MatchedBy(func(u entity.User) bool {
		return u.ID == user.ID
	})).Return("new-access-token-jwt", nil)
	mockRefreshTokenRepo.Mock.On("Save", mock.Anything, mock.MatchedBy(func(rt entity.RefreshToken) bool {
		return rt.UserID == user.ID
	})).Return(errors.New("db save failed"))

	request := httptest.NewRequest(fiber.MethodPost, "/auth/refresh", nil)
	request.Header.Set("Cookie", "refresh_token=db-err-save-token-777")
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

	token1 := "test-grace-t1-" + uuid.New().String()
	mockRefreshTokenRepo.Mock.On("FindByToken", mock.Anything, token1).Return(entity.RefreshToken{
		ID:        "grace-token-id",
		UserID:    user.ID,
		Token:     token1,
		ExpiresAt: time.Now().Add(time.Hour),
	}, nil)
	mockUserRepo.Mock.On("FindByID", mock.Anything, user.ID).Return(user, nil)
	mockJwtUtil.Mock.On("GenerateAccessToken", mock.Anything, mock.MatchedBy(func(u entity.User) bool {
		return u.ID == user.ID
	})).Return("grace-access-token-123", nil)
	mockRefreshTokenRepo.Mock.On("Save", mock.Anything, mock.Anything).Return(nil)

	res1, err := testUsecase.RotateToken(context.Background(), token1)
	assert.NoError(t, err)
	assert.Equal(t, "grace-access-token-123", res1.AccessToken)
	assert.NotEmpty(t, res1.RefreshToken)

	res2, err := testUsecase.RotateToken(context.Background(), token1)
	assert.NoError(t, err)
	assert.Equal(t, res1.AccessToken, res2.AccessToken)
	assert.Equal(t, res1.RefreshToken, res2.RefreshToken)
}
