FROM golang:latest@sha256:f44f6e88636cfb311f9ebace870ded69d943f227bb3cb27d32ffd84ea18c43ea AS build

ARG CGO_ENABLED=0
ARG VERSION
WORKDIR /src
RUN go install github.com/erdnaxeli/wishlister/pkg/cmd@v${VERSION}

FROM cgr.dev/chainguard/static:latest@sha256:bf639cba19ba56329e6907ac26a7afcdde57a80b6aa66d5100da6883196e6b82

WORKDIR /app
COPY --from=build /go/bin/cmd /app/server

CMD ["/app/server"]
