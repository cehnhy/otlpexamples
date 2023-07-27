package trace

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var DefaultHTTPClient = &http.Client{
	Transport: NewTraceTransport(http.DefaultTransport),
}

type TraceTransport struct {
	transport http.RoundTripper
}

func NewTraceTransport(transport http.RoundTripper) *TraceTransport {
	return &TraceTransport{
		transport: transport,
	}
}

func (t *TraceTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	return t.transport.RoundTrip(req)
}
