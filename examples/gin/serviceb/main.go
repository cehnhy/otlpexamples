package main

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/cehnhy/otlpexamples/trace"
)

const (
	_prot           = ":8081"
	_traceName      = ""
	_serviceName    = "serviceB"
	_serviceVersion = "v0.0.1"
)

func main() {
	ctx := context.Background()
	trace.SetPropagator()
	trace.SetTraceProvider(ctx, _serviceName, _serviceVersion)

	e := gin.New()
	e.Use(trace.NewGinMiddleware(_traceName))

	e.GET("/bar", func(c *gin.Context) {
		c.JSON(200, nil)
	})

	e.Run(_prot)
}
