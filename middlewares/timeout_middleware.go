package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"jing-sync/config"
	"jing-sync/logger"
	"jing-sync/utils"
)

func Timeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		done := make(chan struct{})

		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-ctx.Done():
			logger.WebError(c, config.ErrTypeTimeout)
			utils.ResponseError(c, http.StatusGatewayTimeout, config.ErrTypeMsg[config.ErrTypeTimeout])
		case <-done:
		}
	}
}
