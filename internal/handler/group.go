package handler

import (
	"net/http"
	"strings"

	"vincadrn.com/santuy/internal/model"
)

func GetGroupHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		groupId := r.PathValue("groupId")

		if groupId == "" || strings.Contains(groupId, "/") {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		group := model.Group{
			Id:   groupId,
			Name: "group123",
		}

		w.Write([]byte(group.ToJSON()))
	})
}
