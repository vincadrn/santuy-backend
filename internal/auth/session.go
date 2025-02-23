package auth

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"

	"vincadrn.com/santuy/internal/model"
)

const (
	API_SESSION_NAME string = "santuysrv"
)

var cookieStore *sessions.CookieStore

func init() {
	cookieStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	cookieStore.Options = &sessions.Options{
		Path: "/",
		// MaxAge:   7 * 24 * 60 * 60,
		MaxAge:   10 * 60,
		HttpOnly: true,
		Secure:   true,
	}
}

func Session() *sessions.CookieStore {
	return cookieStore
}

// Ideally, all sessions-related operation (read/write)
// should be done here.
// This is to ensure the actual endpoints have received enough information
// stored in the session to do the corresponding CRUDs.
// TODO: Migrate all sessions-related operation here instead.
// * Consider to use this as a wrapper.
func GetSessionValue(key string, w http.ResponseWriter, r *http.Request) any {
	store := Session()
	session, err := store.Get(r, API_SESSION_NAME)
	if err != nil {
		model.ResponseWithErrorDefault(w, err, http.StatusInternalServerError)
		return nil
	}

	return session.Values[key]
}

func SetSessionValue(key string, val any, w http.ResponseWriter, r *http.Request) {
	store := Session()
	session, err := store.Get(r, API_SESSION_NAME)
	if err != nil {
		model.ResponseWithErrorDefault(w, err, http.StatusInternalServerError)
		return
	}

	session.Values[key] = val
	session.Save(r, w)
}
