package model

import (
	"time"
)

type Itinerary struct {
	Date          time.Time `json:"date"`
	Picture       string    `json:"picture"`
	PreparationID string    `json:"preparation_id"`
}

type ItineraryDetail struct {
	ItineraryID string    `json:"itinerary_id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Activity    string    `json:"activity"`
	Picture     string    `json:"picture"`
}

func (i Itinerary) ToJSON() string {
	return MarshalToJSON(i)
}

func (id ItineraryDetail) ToJSON() string {
	return MarshalToJSON(id)
}
