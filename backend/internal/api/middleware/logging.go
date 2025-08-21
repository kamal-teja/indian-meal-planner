package middleware

import (
	"bytes"
	"io"
	"time"

	"nourish-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware creates request logging middleware
func LoggingMiddleware(log *logger.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.Info("HTTP Request",
			"method", param.Method,
			"path", param.Path,
			"status", param.StatusCode,
			"latency", param.Latency,
			"ip", param.ClientIP,
			"user-agent", param.Request.UserAgent(),
			"size", param.BodySize,
		)
		return ""
	})
}

// RequestResponseLogger creates detailed request/response logging middleware
func RequestResponseLogger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Log request
		var reqBody []byte
		if c.Request.Body != nil {
			reqBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		// Capture response
		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:          &bytes.Buffer{},
		}
		c.Writer = w

		// Process request
		c.Next()

		// Log response
		latency := time.Since(start)
		
		log.Info("Request completed",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", latency,
			"ip", c.ClientIP(),
			"request_size", len(reqBody),
			"response_size", w.body.Len(),
		)

		// Log errors if any
		if len(c.Errors) > 0 {
			log.Error("Request errors", "errors", c.Errors.String())
		}
	}
}

// responseWriter captures response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
