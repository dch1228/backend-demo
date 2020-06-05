package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/duchenhao/backend-demo/internal/model"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		ctxI, ok := c.Get("ctx")
		if !ok {
			return
		}

		ctx := ctxI.(*model.ReqContext)

		timeTakenMs := int(time.Since(start) / time.Millisecond)
		status := c.Writer.Status()

		logFields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", status),
			zap.String("remote_addr", c.ClientIP()),
			zap.Int("time_ms", timeTakenMs),
			zap.Int("size", c.Writer.Size()),
		}

		if status >= 500 && status < 600 {
			ctx.Logger.Error("Request Completed", logFields...)
		} else {
			ctx.Logger.Info("Request Completed", logFields...)
		}
	}
}
