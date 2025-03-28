package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// GinContextLogger is a middleware that adds a logger to the context of the request.
func GinContextLogger(log *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		logger := log.With().Logger()
		c.Request = c.Request.WithContext(logger.WithContext(ctx))

		c.Next()
	}
}
