FROM golang:latest@sha256:1ecc479bc712a6bdb56df3e346e33edcc141f469f82840bab9f4bc2bc41bf91d AS build

ARG CGO_ENABLED=0
WORKDIR /src
COPY . .
RUN go build -o /src/server ./cmd/server

FROM cgr.dev/chainguard/static:latest@sha256:a301031ffd4ed67f35ca7fa6cf3dad9937b5fa47d7493955a18d9b4ca5412d1a

WORKDIR /app
COPY --from=build /src/server /app/server

CMD ["/app/server"]
