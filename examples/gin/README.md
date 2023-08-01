# OTLP Gin Examples

## Install

+ install [grafana](https://grafana.com/docs/grafana/latest/setup-grafana/installation/docker/)
+ install [Jaeger all-in-one](https://www.jaegertracing.io/docs/1.47/deployment/#all-in-one)
+ install [loki](https://grafana.com/docs/loki/latest/installation/docker/)
+ install [promtail](https://grafana.com/docs/loki/latest/clients/promtail/installation/) to `./bin/promtail`
+ install [prometheus](https://prometheus.io/docs/prometheus/latest/installation/) with [prometheus-config.yaml](./prometheus-config.yaml)(chage targets)

## Run

```shell
# export OTLP_HTTP_ENDPOINT="localhost:4318"
go run examples/gin/servicea/main.go
go run examples/gin/serviceb/main.go
./bin/promtail -config.file=examples/gin/promtail-config.yaml
curl http://127.0.0.1:8080/foo/0
curl http://127.0.0.1:8080/foo/1
```

## Query

+ Configure Loki  on grafana
+ Configure Jaeger on grafana
+ Configure Prometheus on grafana
+ Add Loki Query On Jaeger: {service="example"} |= `traceId=${__trace.traceId}` |= `spanId=${__span.spanId}`
