FROM golang:latest@sha256:9edf71320ef8a791c4c33ec79f90496d641f306a91fb112d3d060d5c1cee4e20 AS build

ARG CGO_ENABLED=0
ARG VERSION
WORKDIR /src
RUN go install github.com/erdnaxeli/wishlister/pkg/cmd@v${VERSION}

FROM cgr.dev/chainguard/static:latest@sha256:11ec91f0372630a2ca3764cea6325bebb0189a514084463cbb3724e5bb350d14

WORKDIR /app
COPY --from=build /go/bin/cmd /app/server

CMD ["/app/server"]
