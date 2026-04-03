package repository

import (
	"context"
	"database/sql"
	"log"
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
		SELECT i.id, i.date
		FROM travel.itinerary i
		JOIN travel.vacation_vacation_group vvg
			ON i.vacation_id = vvg.vacation_id
		WHERE vvg.group_id = $1
		AND i.vacation_id = $2
		ORDER BY i.date ASC
		`,
		groupId, vacationId,
	)

	if err != nil {
		slog.Error("Cannot list itineraries", "group", groupId, "vacation", vacationId)
		slog.Error(err.Error())

		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatalf("cannot close sql rows iteration: %v", err)
		}
	}()

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
		SELECT il.id, il.itinerary_id, il.start_time, il.end_time, il.activity
		FROM travel.itinerary_list il
		JOIN travel.itinerary i
			ON il.itinerary_id = i.id
		JOIN travel.vacation_vacation_group vvg
			ON i.vacation_id = vvg.vacation_id
		WHERE i.id = $2
		AND vvg.group_id = $1
		ORDER BY il.start_time ASC
		`,
		groupId, itineraryId,
	)

	if err != nil {
		slog.Error("Cannot list itinerary details", "itinerary", itineraryId)
		slog.Error(err.Error())

		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatalf("cannot close sql rows iteration: %v", err)
		}
	}()

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
