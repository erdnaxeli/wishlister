FROM golang:latest@sha256:6cc2338c038bc20f96ab32848da2b5c0641bb9bb5363f2c33e9b7c8838f9a208 AS build

ARG CGO_ENABLED=0
WORKDIR /src
COPY . .
RUN go build -o /src/server ./cmd/server

FROM cgr.dev/chainguard/static:latest@sha256:1ff7590cbc50eaaa917c34b092de0720d307f67d6d795e4f749a0b80a2e95a2c

WORKDIR /app
COPY --from=build /src/server /app/server

CMD ["/app/server"]
