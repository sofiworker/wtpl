package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

var (
	TimeLimiterCount          = 60
	TimeLimiterDuration int64 = 60 // 单位 s
	TokenLimit                = 60
	TokenBurst                = 60
	TokenRate                 = 200 // 单位 ms
)

func TimeLimiterHandler() gin.HandlerFunc {
	start := time.Now().Unix()
	var count int
	return func(c *gin.Context) {
		now := time.Now().Unix()
		dur := now - start
		if dur > TimeLimiterDuration {
			count = 1
			start = now
			c.Next()
			return
		}
		count++
		if count > TimeLimiterCount {
			c.Header("X-Rate-Limit-Limit", fmt.Sprintf("%d", TimeLimiterCount))
			c.Header("X-Rate-Limit-Remaining", "0")
			c.Header("X-Rate-Limit-Reset", fmt.Sprintf("%d", dur))
			c.Status(http.StatusTooManyRequests)
			c.Abort()
			return
		}
		c.Next()
	}
}

func TokenLimiterHandler() gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(time.Duration(TokenRate)*time.Millisecond), TokenBurst)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.Header("X-Rate-Limit-Limit", fmt.Sprintf("%d", TimeLimiterCount))
			c.Header("X-Rate-Limit-Remaining", "0")
			c.Header("X-Rate-Limit-Reset", "0")
			c.Status(http.StatusTooManyRequests)
			c.Abort()
			return
		}
		c.Next()
	}
}
