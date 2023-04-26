package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	"wtpl/pkg/logger"
)

// LoggerHandler instance a Logger middleware with config.
func LoggerHandler(notlogged ...string) gin.HandlerFunc {

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {

			// Stop timer
			timeStamp := time.Now()
			latency := timeStamp.Sub(start)
			if latency > time.Minute {
				latency = latency.Truncate(time.Second)
			}

			clientIP := c.ClientIP()
			method := c.Request.Method
			statusCode := c.Writer.Status()
			//errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

			//bodySize := c.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}

			logger.LDebug("req", zap.String("path", path), zap.String("method", method),
				zap.Int("code", statusCode), zap.Duration("cost", latency), zap.String("client", clientIP))
		}
	}
}
