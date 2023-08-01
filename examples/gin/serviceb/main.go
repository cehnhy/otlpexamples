package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"

	"github.com/cehnhy/otlpexamples/internal/log"
	"github.com/cehnhy/otlpexamples/internal/metric"
	"github.com/cehnhy/otlpexamples/internal/trace"
)

const (
	_prot           = ":8081"
	_traceName      = ""
	_meterName      = "serviceb"
	_serviceName    = "serviceb"
	_serviceVersion = "v0.0.1"
	_logFilePath    = "examples/gin/serviceb.log"
)

func main() {
	ctx := context.Background()
	trace.SetPropagator()
	trace.SetTraceProvider(ctx, _serviceName, _serviceVersion)
	metric.SetMeterProvider(ctx, _serviceName, _serviceVersion)

	log.SetSlog(_logFilePath)

	e := gin.New()
	e.Use(trace.NewGinMiddleware(_traceName))
	e.Use(metric.NewGinMiddleware(_meterName))

	e.GET("/bar", func(c *gin.Context) {
		slog.InfoCtx(c.Request.Context(), "call /bar")
		c.JSON(200, nil)
	})

	e.GET("/metrics", func(c *gin.Context) {
		handler := promhttp.Handler()
		handler.ServeHTTP(c.Writer, c.Request)
	})

	e.Run(_prot)
}
