FROM golang:latest@sha256:b54cbf583d390341599d7bcbc062425c081105cc5ef6d170ced98ef9d047c716 AS build

ARG CGO_ENABLED=0
ARG VERSION
WORKDIR /src
RUN go install github.com/erdnaxeli/wishlister/pkg/cmd@v${VERSION}

FROM cgr.dev/chainguard/static:latest@sha256:bee469c98ce2df388a2746dc360bd59eb3efe2dab366a01cdcbfd738a2ca1474

WORKDIR /app
COPY --from=build /go/bin/cmd /app/server

CMD ["/app/server"]
