package trace

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

func NewGinMiddleware(traceName string) gin.HandlerFunc {
	t := otel.Tracer(traceName)
	return func(c *gin.Context) {
		if c.FullPath() == "/metrics" {
			c.Next()
			return
		}

		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		ctx, span := t.Start(ctx, c.FullPath(), trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		if err := c.Errors.Last(); err != nil {
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, codes.Ok.String())
		}

		// TODO: apply HTTP Server semantic conventions
		span.SetAttributes(semconv.HTTPMethod(c.Request.Method))
		span.SetAttributes(semconv.HTTPURL(c.Request.URL.String()))
		span.SetAttributes(semconv.HTTPScheme(c.Request.URL.Scheme))
		span.SetAttributes(semconv.HTTPStatusCode(c.Writer.Status()))
	}
}
