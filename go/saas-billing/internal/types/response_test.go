package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSuccessResponse(t *testing.T) {
	// Test data
	data := map[string]string{"key": "value"}

	// Test with nil metadata
	resp := NewSuccessResponse(data, nil)
	assert.True(t, resp.Success)
	assert.Equal(t, data, resp.Data)
	assert.NotNil(t, resp.Metadata)
	assert.False(t, resp.Metadata.Timestamp.IsZero())

	// Test with custom metadata
	metadata := &Metadata{
		Timestamp: time.Now(),
	}
	resp = NewSuccessResponse(data, metadata)
	assert.True(t, resp.Success)
	assert.Equal(t, data, resp.Data)
	assert.Equal(t, metadata, resp.Metadata)
}

func TestNewErrorResponse(t *testing.T) {
	// Test data
	errInfo := &ErrorInfo{
		Code:       "INVALID_INPUT",
		Message:    "Invalid input provided",
		Details:    "Field 'email' is required",
		StatusCode: 400,
	}

	resp := NewErrorResponse(errInfo)
	assert.False(t, resp.Success)
	assert.Equal(t, errInfo, resp.Error)
	assert.NotNil(t, resp.Metadata)
	assert.False(t, resp.Metadata.Timestamp.IsZero())
	assert.False(t, resp.Error.Timestamp.IsZero())
}

func TestNewPaginatedResponse(t *testing.T) {
	// Test data
	data := []string{"item1", "item2", "item3"}
	page := 2
	pageSize := 3
	totalRecords := 8

	resp := NewPaginatedResponse(data, page, pageSize, totalRecords)
	assert.True(t, resp.Success)
	assert.Equal(t, data, resp.Data)
	assert.NotNil(t, resp.Metadata)
	assert.NotNil(t, resp.Metadata.Pagination)

	pagination := resp.Metadata.Pagination
	assert.Equal(t, page, pagination.CurrentPage)
	assert.Equal(t, pageSize, pagination.PageSize)
	assert.Equal(t, 3, pagination.TotalPages) // (8 + 2) / 3 = 3
	assert.Equal(t, totalRecords, pagination.TotalRecords)
	assert.True(t, pagination.HasNext)
	assert.True(t, pagination.HasPrevious)
}
