package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("APIError: %s - %s", e.Code, e.Message)
}

func NewAPIError(code string, message string, status int) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

func (e *APIError) Write(w http.ResponseWriter) {
	w.WriteHeader(e.Status)
	json.NewEncoder(w).Encode(e)
}
