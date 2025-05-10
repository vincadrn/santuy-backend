package model

import (
	"time"
)

type Itinerary struct {
	Id   string
	Date time.Time
}

type Itineraries []Itinerary

type ItineraryDetail struct {
	Id          string
	ItineraryId string
	StartTime   TimeOnly
	EndTime     TimeOnly
	Activity    string
}

type ItineraryDetails []ItineraryDetail
