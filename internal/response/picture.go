package response

import (
	"vincadrn.com/santuy/internal/model"
)

type Picture struct {
	Uri string `json:"picture_uri"`
}

type Pictures []Picture

type PicturesWithItinerary struct {
	ItineraryId string `json:"itinerary_id"`
	Pictures    `json:"pictures"`
}

type PicturesWithItineraries []PicturesWithItinerary

func (res Picture) ToJSON() ([]byte, error) {
	return MarshalToJSON(res)
}

func (res Pictures) ToJSON() ([]byte, error) {
	return MarshalToJSON(res)
}

func (res *Picture) Construct(pict *model.Picture) {
	res.Uri = pict.Uri
}

func (res *Pictures) Construct(picts *model.Pictures) {
	for _, item := range *picts {
		var pic Picture
		pic.Construct(&item)

		*res = append(*res, pic)
	}
}

func (res PicturesWithItinerary) ToJSON() ([]byte, error) {
	return MarshalToJSON(res)
}

func (res PicturesWithItineraries) ToJSON() ([]byte, error) {
	return MarshalToJSON(res)
}
