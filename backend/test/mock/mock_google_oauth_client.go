package mock

import (
	"context"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockGoogleOAuthClient struct {
	Mock mock.Mock
}

func (m *MockGoogleOAuthClient) GetLoginURL() string {
	args := m.Mock.Called()
	return args.String(0)
}

func (m *MockGoogleOAuthClient) GetUserInfo(ctx context.Context, code string) (model.GoogleUserInfo, error) {
	args := m.Mock.Called(ctx, code)
	if args.Get(0) == nil {
		return model.GoogleUserInfo{}, args.Error(1)
	}
	return args.Get(0).(model.GoogleUserInfo), args.Error(1)
}
