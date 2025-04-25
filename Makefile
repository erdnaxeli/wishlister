export EMAIL ?= off


all: build-server build-frontend

build-server: generate-repository generate-templates
	go build -o server ./pkg/cmd

build-frontend: build-css

build-css:
	echo "nothing to do"

generate-repository:
	go tool sqlc generate

generate-templates:
	go run statictemplates/cmd/main.go pkg/server/templates/ server pkg/server/

run:
	go tool modd

style:
	go tool golangci-lint fmt ./...
	go tool golangci-lint run ./...

# Release artifact

build-docker-image:
ifndef VERSION
	$(error The variable VERSION must be defined)
endif
	docker build -t ghcr.io/erdnaxeli/wishlister:${VERSION} .

publish: build-docker-image
	docker push ghcr.io/erdnaxeli/wishlister:${VERSION}
