package sso

import (
	"os"

	"github.com/conku/cache/memory"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	GoogleSsoConfig *oauth2.Config
	MemoryCache     = memory.New()
)

func Setup() {
	GoogleSsoConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("SSO_CALLBACK_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"),

		Scopes:   []string{"email", "profile"},
		Endpoint: google.Endpoint,
	}

}
