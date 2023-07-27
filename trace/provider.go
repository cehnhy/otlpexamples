package trace

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func SetTraceProvider(ctx context.Context, serviceName, serviceVersion string) {
	endpoint := os.Getenv("OTLP_HTTP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4318"
	}
	expoter, err := NewExporter(ctx, endpoint)
	if err != nil {
		log.Fatal(err)
	}

	resource := NewResource(serviceName, serviceVersion)

	tp := trace.NewTracerProvider(
		trace.WithBatcher(expoter),
		trace.WithResource(resource),
	)

	otel.SetTracerProvider(tp)
}

func NewExporter(ctx context.Context, endpoint string) (trace.SpanExporter, error) {
	return otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(endpoint), otlptracehttp.WithInsecure())
}

func NewResource(serviceName, version string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion(version),
	)
}
