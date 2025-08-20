package middlewares

import (
	"jing-sync/boot/config"
	"jing-sync/boot/logger"
	"jing-sync/internal/utils"

	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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
