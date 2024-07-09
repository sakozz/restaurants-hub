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
	var config SsoConfig

	switch provider {
	case consts.Google:
		config = SsoConfig{
			oauth2.Config{
				ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
				ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"),
				Endpoint:     google.Endpoint,
				RedirectURL:  os.Getenv("GOOGLE_SSO_CALLBACK_URL"),
				Scopes:       []string{"email", "profile"},
			},
			os.Getenv("GOOGLE_SSO_USER_INFO_URL"),
		}
	default:
		config = SsoConfig{
			oauth2.Config{
				ClientID:     os.Getenv("AUTHENTIK_CLIENT_ID"),
				ClientSecret: os.Getenv("AUTHENTIK_SECRET_KEY"),
				Endpoint:     AuthentikEndpoints(),
				RedirectURL:  os.Getenv("AUTHENTIK_SSO_CALLBACK_URL"),
				Scopes:       []string{"email", "profile", "offline_access", "openid"},
			},
			os.Getenv("AUTHENTIK_SSO_USER_INFO_URL"),
		}
	}

	return &config

}
