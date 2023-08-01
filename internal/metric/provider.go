package metric

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

func SetMeterProvider(ctx context.Context, serviceName, serviceVersion string) {
	expoter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}

	mp := metric.NewMeterProvider(
		metric.WithReader(expoter),
	)

	otel.SetMeterProvider(mp)
}
