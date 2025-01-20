package model

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"
	"time"
)

func TestUserJSON(t *testing.T) {
	user := User{
		Id:    "1",
		Name:  "user123",
		Email: "user@example.com",
		Groups: []GroupRole{
			{
				GroupId: "1",
				Role:    "admin",
			},
			{
				GroupId: "2",
				Role:    "member",
			},
		},
	}

	buffer := new(bytes.Buffer)
	response := []byte(`{
		"id": "1",
		"name": "user123",
		"email": "user@example.com",
		"groups": [
			{
				"group_id": "1",
				"role": "admin"
			},
			{
				"group_id": "2",
				"role": "member"
			}
		]
	}`)

	if err := json.Compact(buffer, response); err != nil {
		log.Fatal(err)
	}

	actual := user.ToJSON()
	expected := buffer.String()

	if actual != expected {
		t.Fatalf(`Want %q, but returned %q`, actual, expected)
	}
}

func TestItineraryJSON(t *testing.T) {
	itinerary := Itinerary{
		Date:          time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
		Picture:       "https://google.com",
		PreparationID: "1",
	}

	actual := itinerary.ToJSON()
	expected := `{"date":"2025-01-02T00:00:00Z","picture":"https://google.com","preparation_id":"1"}`

	if actual != expected {
		t.Fatalf(`Want %q, but returned %q`, actual, expected)
	}
}

func TestItineraryDetailJSON(t *testing.T) {
	id := ItineraryDetail{
		ItineraryID: "1",
		StartTime:   time.Date(2025, 1, 2, 7, 30, 0, 0, time.UTC),
		EndTime:     time.Date(2025, 1, 3, 10, 0, 0, 0, time.UTC),
		Activity:    "activity1",
		Picture:     "https://google.com",
	}

	actual := id.ToJSON()
	expected := `{"itinerary_id":"1","start_time":"2025-01-02T07:30:00Z","end_time":"2025-01-03T10:00:00Z","activity":"activity1","picture":"https://google.com"}`

	if actual != expected {
		t.Fatalf(`Want %q, but returned %q`, actual, expected)
	}
}
