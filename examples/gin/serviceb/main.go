package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"

	"github.com/cehnhy/otlpexamples/internal/log"
	"github.com/cehnhy/otlpexamples/internal/trace"
)

const (
	_prot           = ":8081"
	_traceName      = ""
	_serviceName    = "serviceb"
	_serviceVersion = "v0.0.1"
	_logFilePath    = "examples/gin/serviceb.log"
)

func main() {
	ctx := context.Background()
	trace.SetPropagator()
	trace.SetTraceProvider(ctx, _serviceName, _serviceVersion)

	log.SetSlog(_logFilePath)

	e := gin.New()
	e.Use(trace.NewGinMiddleware(_traceName))

	e.GET("/bar", func(c *gin.Context) {
		slog.InfoCtx(c.Request.Context(), "call /bar")
		c.JSON(200, nil)
	})

	e.Run(_prot)
}
