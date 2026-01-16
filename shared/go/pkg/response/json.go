package response

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// MessageResponse represents a standard message response
type MessageResponse struct {
	Message string `json:"message"`
}

// PaginatedResponse represents a paginated list response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// JSON sends a JSON response with the given status code
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Error sends an error response with the given status code
func Error(w http.ResponseWriter, status int, message, code string) {
	JSON(w, status, ErrorResponse{
		Error: message,
		Code:  code,
	})
}

// ErrorWithDetails sends an error response with details
func ErrorWithDetails(w http.ResponseWriter, status int, message, code, details string) {
	JSON(w, status, ErrorResponse{
		Error:   message,
		Code:    code,
		Details: details,
	})
}

// Success sends a success message response
func Success(w http.ResponseWriter, message string) {
	JSON(w, http.StatusOK, MessageResponse{Message: message})
}

// Created sends a 201 Created response
func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, data)
}

// NoContent sends a 204 No Content response
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Paginated sends a paginated response
func Paginated(w http.ResponseWriter, data interface{}, page, perPage int, total int64) {
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	JSON(w, http.StatusOK, PaginatedResponse{
		Data:       data,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}
