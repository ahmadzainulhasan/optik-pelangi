package mock

import (
	"context"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockGoogleOAuthClient struct {
	Mock mock.Mock
}

func (m *MockGoogleOAuthClient) GetLoginURL(state string, codeChallenge string) string {
	args := m.Mock.Called(state, codeChallenge)
	return args.String(0)
}

func (m *MockGoogleOAuthClient) GetUserInfo(ctx context.Context, code string, codeVerifier string) (model.GoogleUserInfo, error) {
	args := m.Mock.Called(ctx, code, codeVerifier)
	if args.Get(0) == nil {
		return model.GoogleUserInfo{}, args.Error(1)
	}
	return args.Get(0).(model.GoogleUserInfo), args.Error(1)
}
