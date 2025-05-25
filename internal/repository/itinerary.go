package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"vincadrn.com/santuy/internal/model"
)

type ItineraryRepository interface {
	ListItineraries(ctx context.Context, groupId string, vacationId string) (*model.Itineraries, error)
	ListItineraryDetails(ctx context.Context, groupId string, itineraryId string) (*model.ItineraryDetails, error)
}

type itineraryRepository struct {
	db *sql.DB
}

func NewItineraryRepository(db *sql.DB) ItineraryRepository {
	return &itineraryRepository{db: db}
}

func (r *itineraryRepository) ListItineraries(ctx context.Context, groupId string, vacationId string) (*model.Itineraries, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT i.id_itinerary, i.date
		FROM itinerary i
		JOIN groupvacation gv
			ON i.id_vacation = gv.id_vacation
		WHERE gv.id_group = $1
		AND i.id_vacation = $2
		ORDER BY i.date ASC
		`,
		groupId, vacationId,
	)

	if err != nil {
		slog.Error("Cannot list itineraries", "group", groupId, "vacation", vacationId)
		slog.Error(err.Error())

		return nil, err
	}

	defer rows.Close()

	itineraries := &model.Itineraries{}

	for rows.Next() {
		var itinerary model.Itinerary
		err = rows.Scan(&itinerary.Id, &itinerary.Date)
		if err != nil {
			slog.Error("Cannot scan itinerary", "itinerary", itinerary.Id, "date", itinerary.Date)
			slog.Error(err.Error())

			return nil, err
		}
		*itineraries = append(*itineraries, itinerary)
	}

	return itineraries, nil
}

func (r *itineraryRepository) ListItineraryDetails(ctx context.Context, groupId string, itineraryId string) (*model.ItineraryDetails, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT li.id_activities, li.id_itinerary, li.start_time, li.end_time, li.activities
		FROM list_itinerary li
		JOIN itinerary i
			ON li.id_itinerary = i.id_itinerary
		JOIN groupvacation gv
			ON i.id_vacation = gv.id_vacation
		WHERE i.id_itinerary = $2
		AND gv.id_group = $1
		ORDER BY li.start_time ASC
		`,
		groupId, itineraryId,
	)

	if err != nil {
		slog.Error("Cannot list itinerary details", "itinerary", itineraryId)
		slog.Error(err.Error())

		return nil, err
	}

	defer rows.Close()

	itineraryDetails := &model.ItineraryDetails{}
	for rows.Next() {
		var itineraryDetail model.ItineraryDetail
		err = rows.Scan(&itineraryDetail.Id, &itineraryDetail.ItineraryId, &itineraryDetail.StartTime, &itineraryDetail.EndTime, &itineraryDetail.Activity)
		if err != nil {
			slog.Error("Cannot scan itinerary detail", "itinerary_detail", itineraryDetail.Id)
			slog.Error(err.Error())

			return nil, err
		}
		*itineraryDetails = append(*itineraryDetails, itineraryDetail)
	}

	return itineraryDetails, nil
}
