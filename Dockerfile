FROM golang:latest@sha256:ae5a2316d12f3e78fd99177dad452e6ad4f240af2d71d57b480c3477f250fec6 AS build

ARG CGO_ENABLED=0
ARG VERSION
WORKDIR /src
RUN go install github.com/erdnaxeli/wishlister/pkg/cmd@v${VERSION}

FROM cgr.dev/chainguard/static:latest@sha256:60582b2ae6074f641094af0f370d4ab241aab271858a66223dcde7eee9f51638

WORKDIR /app
COPY --from=build /go/bin/cmd /app/server

CMD ["/app/server"]
