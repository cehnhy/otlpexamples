# OTLP Gin Examples

install [Jaeger all-in-one](https://www.jaegertracing.io/docs/1.47/deployment/#all-in-one)

## Run

```shell
# export OTLP_HTTP_ENDPOINT="localhost:4318"
go run examples/gin/servicea/main.go
go run examples/gin/serviceb/main.go
curl http://127.0.0.1:8080/foo/0
```
