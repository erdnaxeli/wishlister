FROM golang:latest@sha256:e2ddb153f786ee6210bf8c40f7f35490b3ff7d38be70d1a0d358ba64225f6428 AS build

ARG CGO_ENABLED=0
ARG VERSION
WORKDIR /src
RUN go install github.com/erdnaxeli/wishlister/pkg/cmd@v${VERSION}

FROM cgr.dev/chainguard/static:latest@sha256:d6a97eb401cbc7c6d48be76ad81d7899b94303580859d396b52b67bc84ea7345

WORKDIR /app
COPY --from=build /go/bin/cmd /app/server

CMD ["/app/server"]
