package model

type GoogleCallbackRequest struct {
	Code        string `query:"code" validate:"required"`
	State       string `query:"state" validate:"required,subtle_eqfield=CookieState"`
	CookieState string `cookie:"oauth_state" validate:"required"`
	Verifier    string `cookie:"pkce_verifier" validate:"required"`
	Error       string `query:"error" validate:"empty"`
}

type GoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type AuthResponse struct {
	AccessToken  string
	RefreshToken string
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type GoogleLoginResult struct {
	URL      string
	State    string
	Verifier string
}
