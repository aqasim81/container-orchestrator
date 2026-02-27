// Package api provides the HTTP API server, router, and response helpers.
package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse is the standard error response body.
type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

// PaginatedResponse wraps a list of items with pagination metadata.
type PaginatedResponse struct {
	Items   any `json:"items"`
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

// JSON writes a JSON response with the given status code.
// It encodes to a buffer first to avoid partial writes on encoding failures.
func JSON(w http.ResponseWriter, status int, data any) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		log.Printf("failed to encode JSON response: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		writeBody(w, []byte(`{"error":"failed to encode response","code":"INTERNAL"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	writeBody(w, buf.Bytes())
}

// Error writes a structured JSON error response.
func Error(w http.ResponseWriter, status int, message, code string) {
	JSON(w, status, ErrorResponse{
		Error: message,
		Code:  code,
	})
}

// Paginated writes a paginated JSON response.
func Paginated(w http.ResponseWriter, items any, total, page, perPage int) {
	JSON(w, http.StatusOK, PaginatedResponse{
		Items:   items,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	})
}

// writeBody writes bytes to the response writer, logging on failure.
func writeBody(w http.ResponseWriter, b []byte) {
	if _, err := w.Write(b); err != nil {
		log.Printf("failed to write response body: %v", err)
	}
}
