FROM golang:latest@sha256:1ecc479bc712a6bdb56df3e346e33edcc141f469f82840bab9f4bc2bc41bf91d AS build

ARG CGO_ENABLED=0
WORKDIR /src
COPY . .
RUN go build -o /src/server ./cmd/server

FROM cgr.dev/chainguard/static:latest@sha256:1ff7590cbc50eaaa917c34b092de0720d307f67d6d795e4f749a0b80a2e95a2c

WORKDIR /app
COPY --from=build /src/server /app/server

CMD ["/app/server"]
