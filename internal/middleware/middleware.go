package middleware

import (
	"log/slog"
	"net/http"

	config "vincadrn.com/santuy/configs"
	"vincadrn.com/santuy/internal/session"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		publicURL := map[string]bool{
			"/v1/auth/login":   true,
			"/v1/auth/session": true,
			"/v1/logout":       true,
			"/oauth2":          true,
		}
		if publicURL[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		sessionProvider := session.NewSessionProvider(w, r)

		_, err := sessionProvider.GetUserEmail()
		if err != nil {
			slog.Error("Cannot retrieve email from session")
			slog.Error(err.Error())
			http.Error(w, "Invalid session", http.StatusBadRequest)

			return
		}

		_, err = sessionProvider.GetUserName()
		if err != nil {
			slog.Error("Cannot retrieve user name from session")
			slog.Error(err.Error())
			http.Error(w, "Invalid session", http.StatusBadRequest)

			return
		}

		next.ServeHTTP(w, r)
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
			_, err := w.Write([]byte("allowed"))
			if err != nil {
				slog.Error("Cannot allow preflight CORS (received OPTIONS method)")
				slog.Error(err.Error())
			}

			return
		}

		next.ServeHTTP(w, r)
	})
}
