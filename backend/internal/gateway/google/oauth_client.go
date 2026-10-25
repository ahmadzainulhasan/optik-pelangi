package google

import (
	"context"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
)

type OAuthClient interface {
	GetLoginURL(state string, codeChallenge string) string
	GetUserInfo(ctx context.Context, code string, codeVerifier string) (model.GoogleUserInfo, error)
}
