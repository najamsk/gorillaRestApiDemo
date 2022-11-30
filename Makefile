ver := najamsk/gorilla_mux_api_8000:1.0.3
lat := najamsk/gorilla_mux_api_8000:latest
img := najamsk/gorilla_mux_api_8000

gver := gorilla_mux_api_8000:1.0.3
glat := gorilla_mux_api_8000:latest

lat := najamsk/gorilla_mux_api_8000:latest
dexe:
	GOOS=linux GOARCH=amd64 go build -o main 1.go

dimg:
	docker build -t ${ver} -t ${lat} .

drun:
	docker run -p 8000:8000 ${lat}

dpush:
	docker image push  ${lat}
	docker image push  ${ver}

gimg:
	docker build -t gcr.io/kubes-369319/${gver} -t gcr.io/kubes-369319/${glat} .

gauth:
	gcloud auth configure-docker

gpush:
	docker image push  gcr.io/kubes-369319/${glat}
	docker image push  gcr.io/kubes-369319/${gver}



.DEFAULT_GOAL := swagger

install_swagger:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger:
	@echo Ensure you have the swagger CLI or this command will fail.
	@echo You can install the swagger CLI with: go get -u github.com/go-swagger/go-swagger/cmd/swagger
	@echo ....

	~/go/bin/swagger generate spec -o ./swagger.json --scan-models
	cp swagger.json swaggerui


swaggery:
	swagger generate spec -o ./swagger.yaml --scan-models
	cp swagger.json swaggerui
