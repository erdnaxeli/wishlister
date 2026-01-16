FROM golang:latest@sha256:bc45dfd319e982dffe4de14428c77defe5b938e29d9bc6edfbc0b9a1fc171cb3 AS build

ARG CGO_ENABLED=0
ARG VERSION
WORKDIR /src
RUN go install github.com/erdnaxeli/wishlister/pkg/cmd@v${VERSION}

FROM cgr.dev/chainguard/static:latest@sha256:a301031ffd4ed67f35ca7fa6cf3dad9937b5fa47d7493955a18d9b4ca5412d1a

WORKDIR /app
COPY --from=build /go/bin/cmd /app/server

CMD ["/app/server"]
