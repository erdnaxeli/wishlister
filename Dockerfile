FROM golang:latest@sha256:d2e5acc5c29cc331ad5e4f59b09dee0f7d47043b072de0799be63e5c49f95feb AS build

ARG CGO_ENABLED=0
ARG VERSION
WORKDIR /src
RUN go install github.com/erdnaxeli/wishlister/pkg/cmd@v${VERSION}

FROM cgr.dev/chainguard/static:latest@sha256:99a5f826e71115aef9f63368120a6aa518323e052297718e9bf084fb84def93c

WORKDIR /app
COPY --from=build /go/bin/cmd /app/server

CMD ["/app/server"]
