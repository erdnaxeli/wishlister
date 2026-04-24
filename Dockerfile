FROM golang:latest@sha256:5f3787b7f902c07c7ec4f3aa91a301a3eda8133aa32661a3b3a3a86ab3a68a36 AS build

ARG CGO_ENABLED=0
ARG VERSION
WORKDIR /src
RUN go install github.com/erdnaxeli/wishlister/pkg/cmd@v${VERSION}

FROM cgr.dev/chainguard/static:latest@sha256:bee469c98ce2df388a2746dc360bd59eb3efe2dab366a01cdcbfd738a2ca1474

WORKDIR /app
COPY --from=build /go/bin/cmd /app/server

CMD ["/app/server"]
