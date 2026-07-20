package response

import (
	"encoding/json"
	"net/http"
)

// Response is the standard API response envelope.
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// JSON writes any value as a JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// Success writes a 200 OK JSON response wrapping the provided data.
func Success(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// Created writes a 201 Created JSON response wrapping the provided data.
func Created(w http.ResponseWriter, data any) {
	JSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// Error writes a JSON error response with the given status code and message.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, Response{
		Success: false,
		Error:   message,
	})
}

// BadRequest writes a 400 Bad Request error response.
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

// NotFound writes a 404 Not Found error response.
func NotFound(w http.ResponseWriter) {
	Error(w, http.StatusNotFound, "resource not found")
}

// InternalError writes a 500 Internal Server Error response.
func InternalError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "internal server error")
}
