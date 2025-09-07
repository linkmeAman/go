package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/saas-billing/internal/logger"
	"github.com/yourusername/saas-billing/internal/types"
)

// Security headers middleware
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		c.Next()
	}
}

// Request ID middleware
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		
		c.Next()
	}
}

// Rate limiting middleware
func RateLimiter() gin.HandlerFunc {
	// Implementation using Redis
	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()
		
		// Check rate limit in Redis
		// ... Redis implementation ...
		
		c.Next()
	}
}

// Request logging middleware
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Process request
		c.Next()
		
		// Log request details
		duration := time.Since(start)
		logger.Info("Request completed", logger.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"duration":   duration.String(),
			"client_ip":  c.ClientIP(),
			"request_id": c.GetString("request_id"),
		})
	}
}

// Error handling middleware
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		
		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// Log error
			logger.Error("Request error", err.Err, logger.Fields{
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"request_id": c.GetString("request_id"),
			})
			
			// Return error response
			c.JSON(http.StatusInternalServerError, types.NewErrorResponse(&types.ErrorInfo{
				Code:      "INTERNAL_SERVER_ERROR",
				Message:   "An unexpected error occurred",
				RequestID: c.GetString("request_id"),
				StatusCode: http.StatusInternalServerError,
			}))
		}
	}
}

// Recovery middleware with custom error handling
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered", nil, logger.Fields{
					"error":      err,
					"request_id": c.GetString("request_id"),
				})
				
				c.JSON(http.StatusInternalServerError, types.NewErrorResponse(&types.ErrorInfo{
					Code:      "INTERNAL_SERVER_ERROR",
					Message:   "An unexpected error occurred",
					RequestID: c.GetString("request_id"),
					StatusCode: http.StatusInternalServerError,
				}))
				
				c.Abort()
			}
		}()
		
		c.Next()
	}
}
