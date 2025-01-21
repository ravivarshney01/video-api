package response

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Response(w http.ResponseWriter, status int, message string, data interface{}, err string) {
	response := jsonResponse{
		Message: message,
		Data:    data,
		Error:   err,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func WithJSON(w http.ResponseWriter, status int, message string, data interface{}) {
	response := jsonResponse{
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func WithError(w http.ResponseWriter, status int, err string) {
	Response(w, status, "", nil, err)
}
