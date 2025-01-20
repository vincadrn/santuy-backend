package handler

import (
	"net/http"
	"strings"

	"vincadrn.com/santuy/internal/model"
)

func GetUserHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("userId")

		// TODO: Add ID/path validation
		if userId == "" || strings.Contains(userId, "/") {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		// for now return a mock user
		user := model.User{
			Id:    userId,
			Name:  "user123",
			Email: "user@example.com",
			Groups: []model.GroupRole{
				{
					GroupId: "1",
					Role:    "admin",
				},
			},
		}

		w.Write([]byte(user.ToJSON()))
	})
}
