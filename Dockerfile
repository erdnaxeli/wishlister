FROM golang:latest@sha256:ce63a16e0f7063787ebb4eb28e72d477b00b4726f79874b3205a965ffd797ab2 AS build

ARG CGO_ENABLED=0
ARG VERSION
WORKDIR /src
RUN go install github.com/erdnaxeli/wishlister/pkg/cmd@v${VERSION}

FROM cgr.dev/chainguard/static:latest@sha256:99a5f826e71115aef9f63368120a6aa518323e052297718e9bf084fb84def93c

WORKDIR /app
COPY --from=build /go/bin/cmd /app/server

CMD ["/app/server"]
