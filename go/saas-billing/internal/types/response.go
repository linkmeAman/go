package types

import "time"

// ApiResponse is the standard response format for all API endpoints
type ApiResponse struct {
	Success  bool        `json:"success"`
	Data     interface{} `json:"data,omitempty"`
	Error    *ErrorInfo  `json:"error,omitempty"`
	Metadata *Metadata   `json:"metadata,omitempty"`
}

// ErrorInfo provides detailed error information
type ErrorInfo struct {
	Code       string    `json:"code"`              // Machine-readable error code
	Message    string    `json:"message"`           // Human-readable error message
	Details    string    `json:"details,omitempty"` // Additional error details
	Timestamp  time.Time `json:"timestamp"`         // When the error occurred
	RequestID  string    `json:"request_id,omitempty"`
	StatusCode int       `json:"-"` // HTTP status code
}

// Metadata includes additional response information like pagination
type Metadata struct {
	Pagination *PaginationInfo `json:"pagination,omitempty"`
	Timestamp  time.Time       `json:"timestamp"`
}

// PaginationInfo provides pagination details
type PaginationInfo struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
	HasNext      bool `json:"has_next"`
	HasPrevious  bool `json:"has_previous"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(data interface{}, metadata *Metadata) *ApiResponse {
	if metadata == nil {
		metadata = &Metadata{
			Timestamp: time.Now(),
		}
	}
	return &ApiResponse{
		Success:  true,
		Data:     data,
		Metadata: metadata,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(err *ErrorInfo) *ApiResponse {
	if err.Timestamp.IsZero() {
		err.Timestamp = time.Now()
	}
	return &ApiResponse{
		Success: false,
		Error:   err,
		Metadata: &Metadata{
			Timestamp: time.Now(),
		},
	}
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, page, pageSize, totalRecords int) *ApiResponse {
	totalPages := (totalRecords + pageSize - 1) / pageSize

	paginationInfo := &PaginationInfo{
		CurrentPage:  page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
		HasNext:      page < totalPages,
		HasPrevious:  page > 1,
	}

	metadata := &Metadata{
		Pagination: paginationInfo,
		Timestamp:  time.Now(),
	}

	return NewSuccessResponse(data, metadata)
}
