package handler

import (
	"context"
	"log/slog"
	"net/http"

	"vincadrn.com/santuy/internal/response"
	"vincadrn.com/santuy/internal/service"
	"vincadrn.com/santuy/internal/session"
)

func UserHandler(svc *service.AccountService, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			slog.Info("Serving GET on userHandler")

			getCurrentUser().ServeHTTP(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
	})
}

func getCurrentUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionProvider := session.NewSessionProvider(w, r)

		email, err := sessionProvider.GetUserEmail()
		if err != nil {
			slog.Error("Cannot retrieve email from session")
			slog.Error(err.Error())
			http.Error(w, "Invalid session", http.StatusBadRequest)

			return
		}

		userName, err := sessionProvider.GetUserName()
		if err != nil {
			slog.Error("Cannot retrieve user name from session")
			slog.Error(err.Error())
			http.Error(w, "Invalid session", http.StatusBadRequest)

			return
		}

		slog.Info("Retrieved username and email", "username", userName, "email", email)

		user := response.User{
			UserName: userName,
			Email:    email,
		}

		response, err := user.ToJSON()
		if err != nil {
			slog.Error("Cannot marshal user", "user", user.UserName)
			slog.Error(err.Error())
			http.Error(w, "Error when serializing user", http.StatusInternalServerError)

			return
		}

		slog.Info("User contains this:", "username", userName, "email", email)
		slog.Info("Sending this response:", "response", response)

		w.Write(response)
	})
}
