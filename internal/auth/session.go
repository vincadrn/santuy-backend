package auth

import (
	"os"

	"github.com/gorilla/sessions"
)

const (
	API_SESSION_NAME string = "santuysrv"
)

var cookieStore *sessions.CookieStore

func init() {
	cookieStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   true,
	}
}

func Session() *sessions.CookieStore {
	return cookieStore
}
