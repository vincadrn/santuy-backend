package response

import (
	"encoding/json"
	"log/slog"
)

type JSONMarshallable interface {
	ToJSON() ([]byte, error)
}

func MarshalToJSON(v JSONMarshallable) ([]byte, error) {
	payload, err := json.Marshal(v)
	if err != nil {
		slog.Error("Cannot marshal", "payload", v)
		slog.Error(err.Error())

		return nil, err
	}

	return payload, nil
}
