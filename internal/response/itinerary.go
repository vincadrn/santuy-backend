package response

import (
	"time"

	"vincadrn.com/santuy/internal/model"
)

type Itinerary struct {
	Id   string    `json:"itinerary_id"`
	Date time.Time `json:"date"`
}

type ItineraryDetail struct {
	Id        string         `json:"detail_id"`
	StartTime model.TimeOnly `json:"start_time"`
	EndTime   model.TimeOnly `json:"end_time"`
	Activity  string         `json:"activity"`
}

type Itineraries []Itinerary

type ItineraryDetails []ItineraryDetail

func (res Itineraries) ToJSON() ([]byte, error) {
	return MarshalToJSON(res)
}

func (res ItineraryDetails) ToJSON() ([]byte, error) {
	return MarshalToJSON(res)
}

func (res *Itinerary) Construct(itinerary *model.Itinerary) {
	res.Id = itinerary.Id
	res.Date = itinerary.Date
}

func (res *Itineraries) Construct(itineraries *model.Itineraries) {
	for _, item := range *itineraries {
		var itinerary Itinerary
		itinerary.Construct(&item)

		*res = append(*res, itinerary)
	}
}

func (res *ItineraryDetail) Construct(itineraryDetails *model.ItineraryDetail) {
	res.Id = itineraryDetails.Id
	res.StartTime = itineraryDetails.StartTime
	res.EndTime = itineraryDetails.EndTime
	res.Activity = itineraryDetails.Activity
}

func (res *ItineraryDetails) Construct(itineraryDetails *model.ItineraryDetails) {
	for _, item := range *itineraryDetails {
		var itineraryDetail ItineraryDetail
		itineraryDetail.Construct(&item)

		*res = append(*res, itineraryDetail)
	}
}
