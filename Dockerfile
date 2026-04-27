FROM golang:latest@sha256:1e598ea5752ae26c093b746fd73c5095af97d6f2d679c43e83e0eac484a33dc3 AS build

ARG CGO_ENABLED=0
ARG VERSION
WORKDIR /src
RUN go install github.com/erdnaxeli/wishlister/pkg/cmd@v${VERSION}

FROM cgr.dev/chainguard/static:latest@sha256:bee469c98ce2df388a2746dc360bd59eb3efe2dab366a01cdcbfd738a2ca1474

WORKDIR /app
COPY --from=build /go/bin/cmd /app/server

CMD ["/app/server"]
