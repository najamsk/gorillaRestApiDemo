# Guide

http server listen on port 8000

## Todos

- [x] Add GET,POST, PUT and DELETE end points and link with swagger
- [x] Embed SwaggerUi folder into output binary
- [x] OpenTelemetry
- [x] Unit tests for http endpoints
- [x] Structured Logging
- [x] Generate http client from swagger
- [ ] Custom error package and struct that can log meaningful trace to zap.logger, currently a detailed jibberish hard read stack is inserted
- [ ] Dockerize
- [ ] Running in Kubernetes inside docker
- [ ] Pushing to AWS?
- [ ] Pushing to GCP?
- [ ] Reading from ENV and map to Config? Viper
- [ ] Reading from Etcd and map to Config? Viper


### Run Jaeger as docker container
```
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.39
  ```
run jaeger ui by visiting http://localhost:16686

Run jaeger container and then run our go server using:

```
go run 1.go
```


### Generating Swagger Go Client

First we will generate swagger client for golang
* on root level create directories gen/client
* run make command make swaggerClient
* once files generated without error run go mod tidy

Now make sure our server is running 

go run 1.go

now cd into cmd/rest/client and run go run client.go
