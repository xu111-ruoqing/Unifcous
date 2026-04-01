package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unifocus/backend/pkg/logger"
)

// MetricsMiddleware 收集API指标
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// 记录指标（简化版本，实际应该发送到Prometheus）
		logger.Infof("[METRICS] %s %s %d %v",
			method,
			path,
			statusCode,
			latency,
		)

		// 设置响应头（可用于前端监控）
		c.Header("X-Response-Time", strconv.FormatInt(latency.Milliseconds(), 10))
		c.Header("X-Request-ID", c.GetString("request_id"))
	}
}
