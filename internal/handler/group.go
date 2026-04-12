package handler

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"vincadrn.com/santuy/internal/model"
	"vincadrn.com/santuy/internal/request"
	"vincadrn.com/santuy/internal/response"
	"vincadrn.com/santuy/internal/service"
	"vincadrn.com/santuy/internal/session"
)

func GroupHandler(svc *service.AccountService, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listGroups(svc, ctx).ServeHTTP(w, r)

		case http.MethodPost:
			joinGroup(svc, ctx).ServeHTTP(w, r)

		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
	})
}

func listGroups(svc *service.AccountService, ctx context.Context) http.Handler {
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

		user := model.User{
			Name:  userName,
			Email: email,
		}

		groups, err := svc.ListGroupsByUser(ctx, &user)
		if err != nil {
			slog.Error("Cannot list group by user")
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		var payload response.GroupRoles
		payload.Construct(groups)

		response, err := payload.ToJSON()
		if err != nil {
			slog.Error("Cannot marshal group roles response")
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		_, err = w.Write(response)
		if err != nil {
			slog.Error("Cannot write response", "response", response)
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}
	})
}

// Only if the group already exists.
// TODO: Create a new function to handle new group creation
func joinGroup(svc *service.AccountService, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request request.JoinGroupRequest
		defer func() {
			err := r.Body.Close()
			if err != nil {
				log.Fatalf("cannot close body in join group: %v", err)
			}
		}()
		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			slog.Error("Cannot read request body when setting group")
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

			return
		}

		sessionProvider := session.NewSessionProvider(w, r)

		email, err := sessionProvider.GetUserEmail()
		if err != nil {
			slog.Error("Cannot retrieve email from session")
			slog.Error(err.Error())
			http.Error(w, "Invalid session", http.StatusBadRequest)

			return
		}

		user, err := svc.GetUserDetails(ctx, email)
		if err != nil {
			slog.Error("Cannot get user from db")
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		requestedGroup := model.Group{
			Id: request.GroupId,
		}
		err = svc.AssignUserToGroup(ctx, user, &requestedGroup)
		if err != nil {
			slog.Error("Cannot assign user to group", "user", email, "group", requestedGroup.Name)
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

			return
		}

		// Write session
		groupRole := session.GroupRole{
			GroupId: requestedGroup.Id,
			Role:    "member",
		}
		err = sessionProvider.SetCurrentGroupRole(groupRole)
		if err != nil {
			slog.Error("Cannot set group ID and role to session", "user", email, "group", requestedGroup.Name)
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		_, err = w.Write([]byte("OK"))
		if err != nil {
			slog.Error("Cannot write response")
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}
	})
}
