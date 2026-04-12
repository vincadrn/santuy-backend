package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"vincadrn.com/santuy/internal/response"
	"vincadrn.com/santuy/internal/service"
	"vincadrn.com/santuy/internal/session"
)

func ItineraryDetailHandler(svc *service.ItineraryService, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listItineraryDetails(svc, ctx).ServeHTTP(w, r)

		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
	})
}

func listItineraryDetails(svc *service.ItineraryService, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		itineraryId := r.PathValue("itineraryId")

		// TODO: Add ID/path validation
		if itineraryId == "" || strings.Contains(itineraryId, "/") {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

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
		slog.Info("Attempting to list itinerary details", "email", email, "groupId", groupRole.GroupId, "role", groupRole.Role)

		itineraryDetails, err := svc.ListItineraryDetails(ctx, groupRole.GroupId, itineraryId)
		if err != nil {
			slog.Error("Cannot list itinerary details", "groupId", groupRole.GroupId, "itineraryId", itineraryId)
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		var itineraryDetailsPayload response.ItineraryDetails
		itineraryDetailsPayload.Construct(itineraryDetails)

		response, err := itineraryDetailsPayload.ToJSON()
		if err != nil {
			slog.Error("Cannot marshal itinerary details to JSON")
			slog.Error(err.Error())
			http.Error(w, "Error when serializing JSON", http.StatusInternalServerError)

			return
		}

		_, err = w.Write(response)
		if err != nil {
			slog.Error("Cannot write response", "response", response)
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	})
}
