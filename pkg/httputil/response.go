package httputil

import (
	"encoding/json"
	"net/http"
)

// ErrorBody is a structured API error response.
type ErrorBody struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// JSON writes a JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Error writes a structured error JSON response.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorBody{
		Error:   http.StatusText(status),
		Message: message,
		Code:    status,
	})
}
