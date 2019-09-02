
run:
	 DOCKER_API_VERSION=$$(docker version --format '{{.Server.APIVersion}}') \
	 go run main.go
test:
	 DOCKER_API_VERSION=$$(docker version --format '{{.Server.APIVersion}}') \
	 go test -cover ./...
build:
	 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -o bin/drwatcher
dockerbuild:
	 docker build -t freundallein/drwatcher:latest .
	 docker push freundallein/drwatcher:latest 