package metric

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"golang.org/x/exp/slog"
)

func NewGinMiddleware(meterName string) gin.HandlerFunc {
	m := otel.Meter(meterName)
	return func(c *gin.Context) {
		if c.FullPath() == "/metrics" {
			c.Next()
			return
		}

		c.Next()

		counter, err := m.Int64Counter("http_request")
		if err != nil {
			slog.Warn("failed to create metric counter", slog.String("error", err.Error()))
			return
		}

		if c.Writer.Status() == http.StatusNotFound {
			opt := metric.WithAttributes(
				attribute.Key("Error").String("NOT_FOUND"),
			)
			counter.Add(c.Request.Context(), 1, opt)
			return
		}

		opt := metric.WithAttributes(
			attribute.Key("Method").String(c.Request.Method),
			attribute.Key("URL").String(c.FullPath()),
		)
		counter.Add(c.Request.Context(), 1, opt)
	}
}
