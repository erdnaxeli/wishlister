FROM golang:latest@sha256:595c7847cff97c9a9e76f015083c481d26078f961c9c8dca3923132f51fe12f1 AS build

ARG CGO_ENABLED=0
ARG VERSION
WORKDIR /src
RUN go install github.com/erdnaxeli/wishlister/pkg/cmd@v${VERSION}

FROM cgr.dev/chainguard/static:latest@sha256:d6d54da1c5bf5d9cecb231786adca86934607763067c8d7d9d22057abe6d5dbc

WORKDIR /app
COPY --from=build /go/bin/cmd /app/server

CMD ["/app/server"]
