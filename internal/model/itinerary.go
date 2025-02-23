package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Itinerary struct {
	Id            string    `json:"id"`
	Date          time.Time `json:"date"`
	Picture       string    `json:"picture"`
	PreparationId string    `json:"preparation_id"`
}

type Itineraries []Itinerary

type ItineraryDetail struct {
	ItineraryId string    `json:"itinerary_id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Activity    string    `json:"activity"`
	Picture     string    `json:"picture"`
}

func (i *Itinerary) GetItinerary(db *sql.DB, ctx context.Context, itineraryId string) {
	err := db.QueryRowContext(
		ctx,
		`SELECT "ID_Itenerary", "Date" FROM "Itenerary" WHERE "ID_Itenerary" = $1;`,
		itineraryId,
	).Scan(&i.Id, &i.Date)

	if err != nil {
		log.Fatal(err)
	}
}

func (i Itinerary) ToJSON() string {
	return MarshalToJSON(i)
}

func (iss Itineraries) ToJSON() string {
	return MarshalToJSON(iss)
}

func (id ItineraryDetail) ToJSON() string {
	return MarshalToJSON(id)
}
