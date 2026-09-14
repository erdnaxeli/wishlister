FROM golang:latest@sha256:f44f6e88636cfb311f9ebace870ded69d943f227bb3cb27d32ffd84ea18c43ea AS build

ARG CGO_ENABLED=0
ARG VERSION
WORKDIR /src
RUN go install github.com/erdnaxeli/wishlister/pkg/cmd@v${VERSION}

FROM cgr.dev/chainguard/static:latest@sha256:207a5673ab31ed83332e54ae33d0f1de4adb5984bd93b8309789889e7bf30ba6

WORKDIR /app
COPY --from=build /go/bin/cmd /app/server

CMD ["/app/server"]
