FROM golang:latest@sha256:ec4debba7b371fb2eaa6169a72fc61ad93b9be6a9ae9da2a010cb81a760d36e7 AS build

ARG CGO_ENABLED=0
ARG VERSION
WORKDIR /src
RUN go install github.com/erdnaxeli/wishlister/pkg/cmd@v${VERSION}

FROM cgr.dev/chainguard/static:latest@sha256:6d508f497fe786ba47d57f4a3cffce12ca05c04e94712ab0356b94a93c4b457f

WORKDIR /app
COPY --from=build /go/bin/cmd /app/server

CMD ["/app/server"]
