package handler

import (
	"net/http"
	"strings"

	"vincadrn.com/santuy/internal/auth"
	"vincadrn.com/santuy/internal/model"
)

func ListGroups() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ggrs := model.GroupGroupRoles{
			model.GroupGroupRole{
				Id:   "1",
				Name: "group1",
				Role: "admin",
			},
			model.GroupGroupRole{
				Id:   "2",
				Name: "group2",
				Role: "member",
			},
		}

		w.Write([]byte(ggrs.ToJSON()))
	})
}

func GroupHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetGroup().ServeHTTP(w, r)
		} else if r.Method == http.MethodPost {
			SetGroup().ServeHTTP(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
	})
}

func GetGroup() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		groupId := r.PathValue("groupId")

		if groupId == "" || strings.Contains(groupId, "/") {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		// return a mock for now
		group := model.Group{
			Id:   groupId,
			Name: "group123",
		}

		w.Write([]byte(group.ToJSON()))
	})
}

func SetGroup() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		groupId := r.PathValue("groupId")

		if groupId == "" || strings.Contains(groupId, "/") {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		// TODO: Add database check if this user is eligible
		// for the requested group and role
		auth.SetSessionValue("group_id", "1", w, r)
		auth.SetSessionValue("role", "Admin", w, r)

		w.Write([]byte("OK"))
	})
}
