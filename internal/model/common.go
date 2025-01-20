package model

import (
	"encoding/json"
	"log"
)

type HttpCode int

type JSONMarshallable interface {
	ToJSON() string
}

func MarshalToJSON(v JSONMarshallable) string {
	payload, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("Cannot marshal to JSON")
	}

	return string(payload)
}
