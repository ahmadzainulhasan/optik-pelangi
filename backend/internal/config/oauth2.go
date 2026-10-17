package config

import (
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewOAuth2(viperConfig *viper.Viper) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     viperConfig.GetString("GOOGLE_CLIENT_ID"),
		ClientSecret: viperConfig.GetString("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  viperConfig.GetString("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

