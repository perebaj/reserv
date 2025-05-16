// Package handler contains all the handlers for the application.
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError is an error that can be written to the response writer
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// NewAPIError creates a new API error
func NewAPIError(code string, message string, status int) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// Error returns the error message
func (e *APIError) Error() string {
	return fmt.Sprintf("APIError: %s - %s", e.Code, e.Message)
}

// Write writes the error to the response writer
func (e *APIError) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(e.Status)

	_, _ = w.Write(b)
}
