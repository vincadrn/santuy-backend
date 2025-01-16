package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"vincadrn.com/santuy/auth"
	"vincadrn.com/santuy/config"
	"vincadrn.com/santuy/model"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		publicURL := map[string]bool{
			"/auth/login": true,
			"/oauth2":     true,
		}
		if publicURL[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		store := auth.Session()
		session, err := store.Get(r, auth.SESSION_NAME)
		if err != nil {
			model.ResponseWithErrorDefault(w, err, http.StatusInternalServerError)
			return
		}

		fooVal := session.Values["foo"]
		log.Println("---- Session in `auth-middleware`:", session.Values)
		if fooVal == "bar" {
			next.ServeHTTP(w, r)
			return
		}

		model.ResponseWithErrorDefault(w, nil, http.StatusForbidden)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := config.Configuration().CORS.AllowedOrigins
		if len(allowedOrigins) > 0 {
			w.Header().Set("Access-Control-Allow-Origin", strings.Join(config.Configuration().CORS.AllowedOrigins, ","))
		}
		if os.Getenv("ENVIRONMENT") == "LOCAL" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT")

		if r.Method == "OPTIONS" {
			w.Write([]byte("allowed"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
