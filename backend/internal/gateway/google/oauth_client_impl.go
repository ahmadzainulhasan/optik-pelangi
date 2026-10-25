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

func (client *OAuthClientImpl) GetLoginURL(state string, codeChallenge string) string {
	return client.oauth2.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
}

func (client *OAuthClientImpl) GetUserInfo(ctx context.Context, code string, codeVerifier string) (model.GoogleUserInfo, error) {
	token, err := client.oauth2.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
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
