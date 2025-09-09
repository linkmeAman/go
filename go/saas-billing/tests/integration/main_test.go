package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/linkmeAman/saas-billing/internal/auth"
	"github.com/linkmeAman/saas-billing/internal/db"
	"github.com/linkmeAman/saas-billing/internal/types"
)

var (
	router *gin.Engine
	token  string
)

func TestMain(m *testing.M) {
	// Set up test environment
	setupTestEnv()
	
	// Run tests
	code := m.Run()
	
	// Cleanup
	cleanupTestEnv()
	
	os.Exit(code)
}

func setupTestEnv() {
	// Set test environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "saas_billing_test")
	os.Setenv("JWT_SECRET", "test_secret")

	// Initialize database
	db.InitDB()

	// Set up router
	router = setupRouter()
}

func cleanupTestEnv() {
	// Clean up database
	db.CloseDB()
}

func TestHealthCheck(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	
	var response types.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

func TestUserRegistration(t *testing.T) {
	// Test data
	userData := map[string]string{
		"email":    "test@example.com",
		"password": "testpassword123",
	}
	
	jsonData, _ := json.Marshal(userData)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	
	var response types.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

func TestUserLogin(t *testing.T) {
	// Test data
	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "testpassword123",
	}
	
	jsonData, _ := json.Marshal(loginData)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	
	var response types.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	
	assert.NoError(t, err)
	assert.True(t, response.Success)
	
	// Store token for other tests
	data := response.Data.(map[string]interface{})
	token = data["token"].(string)
}

func TestCreateOrganization(t *testing.T) {
	// Test data
	orgData := map[string]string{
		"name": "Test Organization",
	}
	
	jsonData, _ := json.Marshal(orgData)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/organizations", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	
	var response types.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	
	assert.NoError(t, err)
	assert.True(t, response.Success)
}
