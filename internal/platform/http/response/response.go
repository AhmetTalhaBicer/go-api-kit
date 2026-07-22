package response

import (
	"encoding/json"
	"net/http"
)

// BaseResponse is the standard response envelope for successful responses with data.
type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorDetail holds the error code and message.
type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"` // For validation error lists and similar payloads.
}

// ErrorResponse is the standard error response envelope.
type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

// JSON writes the payload as JSON to the HTTP client.
func JSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

// Success sends a successful response, such as 200 or 201.
func Success(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	resp := BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	JSON(w, statusCode, resp)
}

// Error sends an error response, such as 400, 404, or 500.
func Error(w http.ResponseWriter, statusCode int, errCode string, message string, details interface{}) {
	resp := ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    errCode,
			Message: message,
			Details: details,
		},
	}
	JSON(w, statusCode, resp)
}
