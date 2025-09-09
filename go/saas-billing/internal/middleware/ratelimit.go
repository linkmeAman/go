package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/linkmeAman/saas-billing/internal/types"
)

func RateLimit(redisClient *redis.Client, key string, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		now := time.Now().UnixNano()
		userKey := fmt.Sprintf("%s:%s", key, c.ClientIP())

		pipe := redisClient.Pipeline()
		pipe.ZRemRangeByScore(ctx, userKey, "0", fmt.Sprint(now-(window.Nanoseconds())))
		pipe.ZAdd(ctx, userKey, &redis.Z{Score: float64(now), Member: now})
		pipe.ZCard(ctx, userKey)
		pipe.Expire(ctx, userKey, window)
		cmds, err := pipe.Exec(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.NewErrorResponse(&types.ErrorInfo{
				Code:       "RATE_LIMIT_ERROR",
				Message:    "Rate limit check failed",
				StatusCode: http.StatusInternalServerError,
			}))
			c.Abort()
			return
		}

		count := cmds[2].(*redis.IntCmd).Val()
		if count > int64(limit) {
			c.JSON(http.StatusTooManyRequests, types.NewErrorResponse(&types.ErrorInfo{
				Code:       "RATE_LIMIT_EXCEEDED",
				Message:    fmt.Sprintf("Rate limit exceeded. Try again in %v", window),
				StatusCode: http.StatusTooManyRequests,
			}))
			c.Abort()
			return
		}

		c.Next()
	}
}

// Example usage in main.go:
// r.Use(middleware.RateLimit(redisClient, "global", 100, time.Minute)) // 100 requests per minute per IP
