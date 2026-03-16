FROM golang:latest@sha256:c7e98cc0fd4dfb71ee7465fee6c9a5f079163307e4bf141b336bb9dae00159a5 AS build

ARG CGO_ENABLED=0
ARG VERSION
WORKDIR /src
RUN go install github.com/erdnaxeli/wishlister/pkg/cmd@v${VERSION}

FROM cgr.dev/chainguard/static:latest@sha256:2fdfacc8d61164aa9e20909dceec7cc28b9feb66580e8e1a65b9f2443c53b61b

WORKDIR /app
COPY --from=build /go/bin/cmd /app/server

CMD ["/app/server"]
