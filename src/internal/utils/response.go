package utils

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Data   interface{} `json:"data"`
	Errors interface{} `json:"errors"`
}

func setDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func SendJson(w http.ResponseWriter, payload interface{}, statusCode int) {
	setDefaultHeaders(w)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&response{
		Data:   payload,
		Errors: nil,
	})
}

func SendJsonError(w http.ResponseWriter, errors interface{}, statusCode int) {
	setDefaultHeaders(w)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&response{
		Data:   nil,
		Errors: errors,
	})
}
