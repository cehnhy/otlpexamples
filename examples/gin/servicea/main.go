package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cehnhy/otlpexamples/trace"
)

const (
	_prot           = ":8080"
	_traceName      = ""
	_serviceName    = "serviceA"
	_serviceVersion = "v0.0.1"
)

const (
	_serviceB = "http://127.0.0.1:8081"
)

func main() {
	ctx := context.Background()
	trace.SetPropagator()
	trace.SetTraceProvider(ctx, _serviceName, _serviceVersion)

	e := gin.New()
	e.Use(trace.NewGinMiddleware(_traceName))

	e.GET("/foo/:id", func(c *gin.Context) {
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

	e.Run(_prot)
}
