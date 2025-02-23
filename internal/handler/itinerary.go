package handler

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"vincadrn.com/santuy/internal/auth"
	"vincadrn.com/santuy/internal/model"
)

func ListItineraries() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, ok := auth.GetSessionValue("email", w, r).(string)
		if !ok {
			http.Error(w, "Invalid session", http.StatusBadRequest)
			return
		}
		groupId, ok := auth.GetSessionValue("group_id", w, r).(string)
		if !ok {
			http.Error(w, "Invalid session", http.StatusBadRequest)
			return
		}
		role, ok := auth.GetSessionValue("role", w, r).(string)
		if !ok {
			http.Error(w, "Invalid session", http.StatusBadRequest)
			return
		}

		// TODO: Add database query to check for
		// groupId and role
		// and then dump all the itineraries
		log.Println("Email:", email, "GroupID:", groupId, "Role:", role)

		// for now return a mock itinerary
		itineraries := model.Itineraries{
			model.Itinerary{
				Id:            email,
				Date:          time.Now(),
				Picture:       "https://google.com",
				PreparationId: "2",
			},
			model.Itinerary{
				Id:            email,
				Date:          time.Now(),
				Picture:       "https://google.com",
				PreparationId: "2",
			},
			model.Itinerary{
				Id:            email,
				Date:          time.Now(),
				Picture:       "https://google.com",
				PreparationId: "2",
			},
		}

		w.Write([]byte(itineraries.ToJSON()))
	})
}

func GetItinerary(db *sql.DB, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		itineraryId := r.PathValue("itineraryId")

		// TODO: Add ID/path validation
		if itineraryId == "" || strings.Contains(itineraryId, "/") {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		itinerary := model.Itinerary{}
		itinerary.GetItinerary(db, ctx, itineraryId)

		w.Write([]byte(itinerary.ToJSON()))
	})
}
