package models

import "encoding/json"

// Money represents a monetary amount.
type Money struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

// PageDetail contains pagination metadata from API responses.
type PageDetail struct {
	TotalResults int `json:"totalResults"`
	StartIndex   int `json:"startIndex"`
	ItemsPerPage int `json:"itemsPerPage"`
}

// APIResponse is the standard wrapper for all API responses.
type APIResponse struct {
	Data       json.RawMessage `json:"data"`
	Pagination *PageDetail     `json:"pagination,omitempty"`
	Error      *ErrorBody      `json:"error,omitempty"`
}

// ErrorBody represents an API error response.
type ErrorBody struct {
	Errors []APIError `json:"errors"`
}

// APIError is a single error from the API.
type APIError struct {
	MessageCode string `json:"messageCode"`
	Message     string `json:"message"`
	Field       string `json:"field,omitempty"`
}
