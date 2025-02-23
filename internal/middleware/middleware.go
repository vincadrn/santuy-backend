package middleware

import (
	"log"
	"net/http"

	config "vincadrn.com/santuy/configs"
	"vincadrn.com/santuy/internal/auth"
	"vincadrn.com/santuy/internal/model"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		publicURL := map[string]bool{
			"/v1/auth/login":   true,
			"/v1/auth/session": true,
			"/oauth2":          true,
		}
		if publicURL[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		store := auth.Session()
		session, err := store.Get(r, auth.API_SESSION_NAME)
		if err != nil {
			model.ResponseWithErrorDefault(w, err, http.StatusInternalServerError)
			return
		}

		log.Println("In auth middleware. Session values:", session)

		if session.Values["email"] == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		userEmail := session.Values["email"].(string)
		if userEmail != "" {
			next.ServeHTTP(w, r)
			return
		}

		model.ResponseWithErrorDefault(w, nil, http.StatusForbidden)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := make(map[string]bool)
		for _, allowedOrigin := range config.Configuration().CORS.AllowedOrigins {
			allowedOrigins[allowedOrigin] = true
		}
		origin := r.Header.Get("Origin")
		if !allowedOrigins[origin] {
			http.Error(w, "Origin not allowed", http.StatusForbidden)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		authEndpoints := map[string]bool{
			"/v1/auth/login":    true,
			"/v1/auth/callback": true,
		}
		if authEndpoints[r.URL.Path] {
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")
		} else {
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT")
		}

		if r.Method == "OPTIONS" {
			w.Write([]byte("allowed"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
