package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func RequestTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		fullPath := c.FullPath()

		tracer := opentracing.GlobalTracer()
		span := tracer.StartSpan(fmt.Sprintf("full_path %s", fullPath))
		defer span.Finish()

		ctx := opentracing.ContextWithSpan(c.Request.Context(), span)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		status := c.Writer.Status()

		ext.HTTPStatusCode.Set(span, uint16(status))
		ext.HTTPMethod.Set(span, c.Request.Method)
		if status >= 400 {
			ext.Error.Set(span, true)
		}
	}
}
