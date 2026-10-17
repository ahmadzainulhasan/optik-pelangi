package google

import (
	"context"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
)

type OAuthClient interface {
	GetLoginURL() string
	GetUserInfo(ctx context.Context, code string) (model.GoogleUserInfo, error)
}
