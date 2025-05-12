package auth

import (
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	config "vincadrn.com/santuy/configs"
)

var (
	conf = config.Configuration()

	OauthConfig = oauth2.Config{}
)

func init() {
	var obtainer OAuthObtainer

	switch conf.SecretKeeper {
	case ORACLE:
		obtainer = OracleOAuthCredentials{}
	case NORTHFLANK:
		obtainer = NorthflankOAuthCredentials{}
	case LOCAL:
		obtainer = LocalOAuthCredentials{}
	}

	OauthConfig = oauth2.Config{
		ClientID:     obtainer.GetOAuthClientID(),
		ClientSecret: obtainer.GetOAuthClientSecret(),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func SetOAuthRedirectURL(host string) {
	OauthConfig.RedirectURL = fmt.Sprintf("https://%s/oauth2", host)
}

type OAuthRedirectURI struct {
	RedirectURI string `json:"redirect_uri"`
}
