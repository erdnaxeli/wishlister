FROM golang:latest@sha256:3680233e3204827fbdc66088528ae6d4b3d034f51d03a99d454f6de034888244 AS build

ARG CGO_ENABLED=0
ARG VERSION
WORKDIR /src
RUN go install github.com/erdnaxeli/wishlister/pkg/cmd@v${VERSION}

FROM cgr.dev/chainguard/static:latest@sha256:8658e9eaeb2294fa0e49b0667f7d66b4e21fa7690b8d2a4c30cd1d013ceb9eb7

WORKDIR /app
COPY --from=build /go/bin/cmd /app/server

CMD ["/app/server"]
