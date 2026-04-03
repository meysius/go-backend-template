package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-starter-template/helpers"
)

const RequestIDHeader = "X-Request-ID"

func RequestLogger(base *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Header(RequestIDHeader, requestID)

		logger := base.With("request_id", requestID)
		ctx := helpers.WithLogger(c.Request.Context(), logger)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		logger.Info("request completed",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", time.Since(start).String(),
			"client_ip", c.ClientIP(),
		)
	}
}
