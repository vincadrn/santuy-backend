package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"vincadrn.com/santuy/internal/request"
	"vincadrn.com/santuy/internal/response"
	"vincadrn.com/santuy/internal/service"
	"vincadrn.com/santuy/internal/session"
)

func ListAllPictures(pictSvc *service.PictureService, itinerarySvc *service.ItineraryService, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
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
		slog.Info("Attempting to list itineraries", "email", email, "groupId", groupRole.GroupId, "role", groupRole.Role)

		// TODO: Make vacationId not hardcoded :(
		vacationId := 1
		itineraries, err := itinerarySvc.ListItineraries(ctx, groupRole.GroupId, strconv.Itoa(vacationId))
		if err != nil {
			slog.Error("Cannot list itineraries", "groupId", groupRole.GroupId, "vacationId", vacationId)
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		payload := response.PicturesWithItineraries{}
		for _, itinerary := range *itineraries {
			pictures, err := pictSvc.ListPicturesByItinerary(ctx, itinerary.Id)
			if err != nil {
				slog.Error("Cannot list pictures by itinerary", "itinerary", itinerary)
				slog.Error(err.Error())

				return
			}

			slog.Info("Got pictures", "pictures", pictures, "itinerary", itinerary)

			var pictsRes response.Pictures
			pictsRes.Construct(pictures)

			payload = append(payload, response.PicturesWithItinerary{
				ItineraryId: itinerary.Id,
				Pictures:    pictsRes,
			})
		}

		res, err := payload.ToJSON()
		if err != nil {
			slog.Error("Cannot marshal pictures with itineraries to JSON")
			slog.Error(err.Error())
			http.Error(w, "Error when serializing JSON", http.StatusInternalServerError)

			return
		}

		_, err = w.Write(res)
		if err != nil {
			slog.Error("Cannot write response", "response", res)
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}
	})
}

func ListPicturesByItinerary(pictSvc *service.PictureService, itinerarySvc *service.ItineraryService, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
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
		slog.Info("Attempting to list itineraries", "email", email, "groupId", groupRole.GroupId, "role", groupRole.Role)

		// TODO: Make vacationId not hardcoded :(
		vacationId := 1
		itineraries, err := itinerarySvc.ListItineraries(ctx, groupRole.GroupId, strconv.Itoa(vacationId))
		if err != nil {
			slog.Error("Cannot list itineraries", "groupId", groupRole.GroupId, "vacationId", vacationId)
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		// TODO: Add ID/path validation
		itineraryId := r.PathValue("itineraryId")
		if itineraryId == "" || strings.Contains(itineraryId, "/") {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		// Validate if itineraryId belongs to the user
		itineraryBelongsToUser := false
		for _, listedItinerary := range *itineraries {
			if listedItinerary.Id == itineraryId {
				itineraryBelongsToUser = true
				break
			}
		}
		if !itineraryBelongsToUser {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		pictures, err := pictSvc.ListPicturesByItinerary(ctx, itineraryId)
		if err != nil {
			slog.Error("Cannot list pictures by itinerary", "itinerary", itineraryId)
			slog.Error(err.Error())

			return
		}

		var payload response.Pictures
		payload.Construct(pictures)

		res, err := payload.ToJSON()
		if err != nil {
			slog.Error("Cannot marshal pictures to JSON")
			slog.Error(err.Error())
			http.Error(w, "Error when serializing JSON", http.StatusInternalServerError)

			return
		}

		_, err = w.Write(res)
		if err != nil {
			slog.Error("Cannot write response", "response", res)
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}
	})
}

// Upload by means of PAR
// so backend does not handle uploading directly
func CreateUploadObjectURI(pictSvc *service.PictureService, itinerarySvc *service.ItineraryService, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		var request request.PictureRequest
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&request)

		// Check if format is not empty
		if request.Format == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
		slog.Info("Attempting to list itineraries", "email", email, "groupId", groupRole.GroupId, "role", groupRole.Role)

		// TODO: Make vacationId not hardcoded :(
		vacationId := 1
		itineraries, err := itinerarySvc.ListItineraries(ctx, groupRole.GroupId, strconv.Itoa(vacationId))
		if err != nil {
			slog.Error("Cannot list itineraries", "groupId", groupRole.GroupId, "vacationId", vacationId)
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		// TODO: Add ID/path validation
		itineraryId := r.PathValue("itineraryId")
		if itineraryId == "" || strings.Contains(itineraryId, "/") {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		// Validate if itineraryId belongs to the user
		itineraryBelongsToUser := false
		for _, listedItinerary := range *itineraries {
			if listedItinerary.Id == itineraryId {
				itineraryBelongsToUser = true
				break
			}
		}
		if !itineraryBelongsToUser {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		picture, err := pictSvc.CreateObjectUploadURI(ctx, itineraryId, request.Format)
		if err != nil {
			slog.Error("Cannot create upload URI", "itinerary", itineraryId)
			slog.Error(err.Error())

			return
		}

		var payload response.Picture
		payload.Construct(picture)

		res, err := payload.ToJSON()
		if err != nil {
			slog.Error("Cannot marshal picture to JSON")
			slog.Error(err.Error())
			http.Error(w, "Error when serializing JSON", http.StatusInternalServerError)

			return
		}

		_, err = w.Write(res)
		if err != nil {
			slog.Error("Cannot write response", "response", res)
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}
	})
}
