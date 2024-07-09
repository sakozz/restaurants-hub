package configs

import (
	"os"

	"github.com/conku/cache/memory"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	consts "resturants-hub.com/m/v2/packages/const"
)

var (
	MemoryCache = memory.New()
)

type SsoConfig struct {
	oauth2.Config
	UserInfoUrl string
}

func AuthentikEndpoints() oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:   os.Getenv("AUTH_URL"),
		TokenURL:  os.Getenv("TOKEN_URL"),
		AuthStyle: oauth2.AuthStyleInParams,
	}
}

func NewSsoConfig(provider consts.SsoProvider) *SsoConfig {

	config := SsoConfig{
		oauth2.Config{
			RedirectURL: os.Getenv("SSO_CALLBACK_URL"),
			Scopes:      []string{"email", "profile", "offline_access", "openid"},
		},
		"",
	}

	switch provider {
	case consts.Google:
		config.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
		config.ClientSecret = os.Getenv("GOOGLE_SECRET_KEY")
		config.UserInfoUrl = os.Getenv("GOOGLE_SSO_USER_INFO_URL")
		config.Endpoint = google.Endpoint
	default:
		config.ClientID = os.Getenv("AUTHENTIK_CLIENT_ID")
		config.ClientSecret = os.Getenv("AUTHENTIK_SECRET_KEY")
		config.UserInfoUrl = os.Getenv("AUTHENTIK_SSO_USER_INFO_URL")
		config.Endpoint = AuthentikEndpoints()
	}

	return &config
}
