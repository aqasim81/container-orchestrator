// Package handlers contains HTTP request handlers for the orchestrator API.
package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// HealthResponse is the response body for the health check endpoint.
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// Health returns an HTTP handler for the health check endpoint.
func Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := HealthResponse{
			Status:    "ok",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(resp); err != nil {
			log.Printf("failed to encode health response: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if _, writeErr := w.Write([]byte(`{"error":"encode failed","code":"INTERNAL"}`)); writeErr != nil {
				log.Printf("failed to write error response: %v", writeErr)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(buf.Bytes()); err != nil {
			log.Printf("failed to write health response: %v", err)
		}
	}
}
