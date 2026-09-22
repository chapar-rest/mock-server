package restutils

import (
	"encoding/json"
	"net/http"
)

// PlainText writes a plain text response.
func PlainText(response http.ResponseWriter, status int, body string) {
	response.Header().Set("Content-Type", "text/plain")
	response.WriteHeader(status)
	_, _ = response.Write([]byte(body))
}

// JSON writes a JSON response.
func JSON(response http.ResponseWriter, status int, data any) {
	body, err := json.Marshal(data)
	if err != nil {
		PlainText(response, http.StatusInternalServerError, "Something went wrong")
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	_, _ = response.Write(body)
}

// Error writes an error response.
func Error(response http.ResponseWriter, status int, err error) {
	JSON(response, status, map[string]string{"error": err.Error()})
}

// Success writes a success response.
func Success(response http.ResponseWriter, status int, data any) {
	JSON(response, status, data)
}
