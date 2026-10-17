package google

import (
	"context"
	"encoding/json"

	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/model"
	"golang.org/x/oauth2"
)

type OAuthClientImpl struct {
	oauth2 *oauth2.Config
}

func NewOAuthClientImpl(oauth2 *oauth2.Config) *OAuthClientImpl {
	return &OAuthClientImpl{
		oauth2: oauth2,
	}
}

func (client *OAuthClientImpl) GetLoginURL() string {
	return client.oauth2.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func (client *OAuthClientImpl) GetUserInfo(ctx context.Context, code string) (model.GoogleUserInfo, error) {
	token, err := client.oauth2.Exchange(ctx, code)
	if err != nil {
		return model.GoogleUserInfo{}, err
	}

	httpClient := client.oauth2.Client(ctx, token)
	response, err := httpClient.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return model.GoogleUserInfo{}, err
	}
	defer response.Body.Close()

	var userInfo model.GoogleUserInfo
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		return model.GoogleUserInfo{}, err
	}

	return userInfo, nil
}
