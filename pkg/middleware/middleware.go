package middlewares

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
)

func LogRequest(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method
		tenantID := c.GetHeader("X-Tenant-ID")

		c.Next()
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		log.Infof(i18n.Translate(ctx, "Method=%s Path=%s Query=%s Status=%d Latency=%s TenantID=%s"),
			method, path, raw, statusCode, latency, tenantID)
	}
}
