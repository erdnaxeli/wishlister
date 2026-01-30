FROM golang:latest@sha256:bc45dfd319e982dffe4de14428c77defe5b938e29d9bc6edfbc0b9a1fc171cb3 AS build

ARG CGO_ENABLED=0
ARG VERSION
WORKDIR /src
RUN go install github.com/erdnaxeli/wishlister/pkg/cmd@v${VERSION}

FROM cgr.dev/chainguard/static:latest@sha256:99a5f826e71115aef9f63368120a6aa518323e052297718e9bf084fb84def93c

WORKDIR /app
COPY --from=build /go/bin/cmd /app/server

CMD ["/app/server"]
