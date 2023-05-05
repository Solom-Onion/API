package util

import (
	"encoding/json"
	"net/http"
)

func SendErrorResponse(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	errorResponse := struct {
		Error string `json:"error"`
	}{
		Error: message,
	}
	json.NewEncoder(w).Encode(errorResponse)
}
