FROM golang:latest@sha256:c83e68f3ebb6943a2904fa66348867d108119890a2c6a2e6f07b38d0eb6c25c5 AS build

ARG CGO_ENABLED=0
ARG VERSION
WORKDIR /src
RUN go install github.com/erdnaxeli/wishlister/pkg/cmd@v${VERSION}

FROM cgr.dev/chainguard/static:latest@sha256:11ec91f0372630a2ca3764cea6325bebb0189a514084463cbb3724e5bb350d14

WORKDIR /app
COPY --from=build /go/bin/cmd /app/server

CMD ["/app/server"]
