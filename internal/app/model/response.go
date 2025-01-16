package model

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Status int    `json:"status"`
	Reason string `json:"reason"`
}

func ResponseWithErrorDefault(w http.ResponseWriter, err error, httpStatus int) {
	ResponseWithError(w, err, httpStatus, "")
}

func ResponseWithError(w http.ResponseWriter, err error, httpStatus int, reason string) {
	if err != nil {
		log.Println(err.Error())
	}

	errorReason := reason
	if errorReason == "" {
		errorReason = http.StatusText(httpStatus)
	}

	resp := ErrorResponse{
		Status: httpStatus,
		Reason: reason,
	}

	payload, e := json.Marshal(resp)
	if e != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	http.Error(w, string(payload), httpStatus)
}
