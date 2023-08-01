package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"

	"github.com/cehnhy/otlpexamples/internal/log"
	"github.com/cehnhy/otlpexamples/internal/metric"
	"github.com/cehnhy/otlpexamples/internal/trace"
)

const (
	_prot           = ":8080"
	_traceName      = ""
	_meterName      = "servicea"
	_serviceName    = "servicea"
	_serviceVersion = "v0.0.1"
	_logFilePath    = "examples/gin/servicea.log"
)

const (
	_serviceB = "http://127.0.0.1:8081"
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

	e.GET("/foo/:id", func(c *gin.Context) {
		slog.InfoCtx(c.Request.Context(), "call /foo/:id", slog.String("id", c.Param("id")))
		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, _serviceB+"/bar", nil)
		if err != nil {
			c.Error(err)
			return
		}

		_, err = trace.DefaultHTTPClient.Do(req)
		if err != nil {
			c.Error(err)
			return
		}

		if c.Param("id") == "0" {
			c.Error(errors.New("id should not be 0"))
		}

		c.JSON(200, nil)
	})

	e.GET("/metrics", func(c *gin.Context) {
		handler := promhttp.Handler()
		handler.ServeHTTP(c.Writer, c.Request)
	})

	e.Run(_prot)
}
