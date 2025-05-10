package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"vincadrn.com/santuy/internal/response"
	"vincadrn.com/santuy/internal/service"
	"vincadrn.com/santuy/internal/session"
)

// func Itinerary(db *sql.DB, ctx context.Context) http.Handler {
func ItineraryHandler(svc *service.ItineraryService, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			listItineraries(svc, ctx).ServeHTTP(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
	})
}

func listItineraries(svc *service.ItineraryService, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionProvider := session.NewSessionProvider(w, r)

		email, err := sessionProvider.GetUserEmail()
		if err != nil {
			slog.Error("Cannot get email from session")
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

			return
		}

		groupRole, err := sessionProvider.GetCurrentGroupRole()
		if err != nil {
			slog.Error("Cannot get current group ID and role from session")
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

			return
		}
		slog.Info("Attempting to list itineraries", "email", email, "groupId", groupRole.GroupId, "role", groupRole.Role)

		// TODO: Make vacationId not hardcoded :(
		vacationId := 1
		itineraries, err := svc.ListItineraries(ctx, groupRole.GroupId, strconv.Itoa(vacationId))
		if err != nil {
			slog.Error("Cannot list itineraries", "groupId", groupRole.GroupId, "vacationId", vacationId)
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		var itinerariesPayload response.Itineraries
		itinerariesPayload.Construct(itineraries)

		response, err := itinerariesPayload.ToJSON()
		if err != nil {
			slog.Error("Cannot marshal itineraries to JSON")
			slog.Error(err.Error())
			http.Error(w, "Error when serializing JSON", http.StatusInternalServerError)

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
