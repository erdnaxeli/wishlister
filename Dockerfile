FROM golang:latest@sha256:6cc2338c038bc20f96ab32848da2b5c0641bb9bb5363f2c33e9b7c8838f9a208 AS build

ARG CGO_ENABLED=0
WORKDIR /src
COPY . .
RUN go build -o /src/server ./cmd/server

FROM cgr.dev/chainguard/static:latest@sha256:a301031ffd4ed67f35ca7fa6cf3dad9937b5fa47d7493955a18d9b4ca5412d1a

WORKDIR /app
COPY --from=build /src/server /app/server

CMD ["/app/server"]
