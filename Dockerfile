FROM golang:latest@sha256:f96cc555eb8db430159a3aa6797cd5bae561945b7b0fe7d0e284c63a3b291609 AS build

ARG CGO_ENABLED=0
ARG VERSION
WORKDIR /src
RUN go install github.com/erdnaxeli/wishlister/pkg/cmd@v${VERSION}

FROM cgr.dev/chainguard/static:latest@sha256:60582b2ae6074f641094af0f370d4ab241aab271858a66223dcde7eee9f51638

WORKDIR /app
COPY --from=build /go/bin/cmd /app/server

CMD ["/app/server"]
