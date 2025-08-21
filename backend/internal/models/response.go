package models

// Common response structures

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"totalPages"`
	HasNext    bool `json:"hasNext"`
	HasPrev    bool `json:"hasPrev"`
}

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
	Code    string      `json:"code,omitempty"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Token   string       `json:"token"`
	User    UserResponse `json:"user"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
	Database  string `json:"database"`
}
