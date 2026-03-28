package response

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"vincadrn.com/santuy/internal/model"
	"vincadrn.com/santuy/internal/response"
)

func TestCreateItinerariesResponseFromSingleItinerary(t *testing.T) {
	nowTime := time.Now()
	itinerariesModel := model.Itineraries{
		model.Itinerary{
			Id:   "1",
			Date: nowTime,
		},
	}

	resp := response.Itineraries{}
	resp.Construct(&itinerariesModel)
	_, err := resp.ToJSON()

	assert.Equal(t, "1", resp[0].Id, "id should be equal")
	assert.Equal(t, nowTime, resp[0].Date, "date should be equal")
	assert.Equal(t, 1, len(resp), "length of itineraries should be equal")
	assert.NoError(t, err, "json marshalling should not produce error")
}

func TestCreateItinerariesResponseFromMultipleItinerary(t *testing.T) {
	nowTime := time.Now()
	itinerariesModel := model.Itineraries{
		model.Itinerary{
			Id:   "1",
			Date: nowTime,
		},
		model.Itinerary{
			Id:   "2",
			Date: nowTime,
		},
		model.Itinerary{
			Id:   "3",
			Date: nowTime,
		},
	}

	resp := response.Itineraries{}
	resp.Construct(&itinerariesModel)
	_, err := resp.ToJSON()

	assert.Equal(t, "2", resp[1].Id, "id should be equal")
	assert.Equal(t, nowTime, resp[1].Date, "date should be equal")
	assert.Equal(t, 3, len(resp), "length of itineraries should be equal")
	assert.NoError(t, err, "json marshalling should not produce error")
}

func TestCreateItineraryDetailsResponseFromSingleItineraryDetail(t *testing.T) {
	nowTime := time.Now()
	itinerariesModel := model.ItineraryDetails{
		model.ItineraryDetail{
			Id:          "1",
			ItineraryId: "1",
			StartTime:   model.TimeOnly{Time: nowTime},
			EndTime:     model.TimeOnly{Time: nowTime},
		},
	}

	resp := response.ItineraryDetails{}
	resp.Construct(&itinerariesModel)
	_, err := resp.ToJSON()

	assert.Equal(t, "1", resp[0].Id, "id should be equal")
	assert.Equal(t, model.TimeOnly{Time: nowTime}, resp[0].StartTime, "start date should be equal")
	assert.Equal(t, model.TimeOnly{Time: nowTime}, resp[0].EndTime, "end date should be equal")
	assert.Equal(t, 1, len(resp), "length of itineraries should be equal")
	assert.NoError(t, err, "json marshalling should not produce error")
}

func TestCreateItineraryDetailsResponseFromMultipleItineraryDetails(t *testing.T) {
	nowTime := time.Now()
	itinerariesModel := model.ItineraryDetails{
		model.ItineraryDetail{
			Id:          "1",
			ItineraryId: "1",
			StartTime:   model.TimeOnly{Time: nowTime},
			EndTime:     model.TimeOnly{Time: nowTime},
		},
		model.ItineraryDetail{
			Id:          "2",
			ItineraryId: "1",
			StartTime:   model.TimeOnly{Time: nowTime},
			EndTime:     model.TimeOnly{Time: nowTime},
		},
		model.ItineraryDetail{
			Id:          "3",
			ItineraryId: "1",
			StartTime:   model.TimeOnly{Time: nowTime},
			EndTime:     model.TimeOnly{Time: nowTime},
		},
	}

	resp := response.ItineraryDetails{}
	resp.Construct(&itinerariesModel)
	_, err := resp.ToJSON()

	assert.Equal(t, "1", resp[0].Id, "id should be equal")
	assert.Equal(t, model.TimeOnly{Time: nowTime}, resp[0].StartTime, "start date should be equal")
	assert.Equal(t, model.TimeOnly{Time: nowTime}, resp[0].EndTime, "end date should be equal")
	assert.Equal(t, 3, len(resp), "length of itineraries should be equal")
	assert.NoError(t, err, "json marshalling should not produce error")
}
