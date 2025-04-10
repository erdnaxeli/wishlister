all: build-server build-frontend

build-server: generate-repository
	go build ./cmd/server

build-frontend: build-css

build-css:
	echo "nothing to do"


generate-repository:
	go tool sqlc generate

run:
	go tool modd

style:
	go tool golangci-lint fmt ./...
	go tool golangci-lint run ./...
