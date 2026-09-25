FROM golang:latest@sha256:f44f6e88636cfb311f9ebace870ded69d943f227bb3cb27d32ffd84ea18c43ea AS build

ARG CGO_ENABLED=0
ARG VERSION
WORKDIR /src
RUN go install github.com/erdnaxeli/wishlister/pkg/cmd@v${VERSION}

FROM cgr.dev/chainguard/static:latest@sha256:41e17ed83c594a64a9396b6ab96dd26d5ddc290dacf4c177464712ff21ad534f

WORKDIR /app
COPY --from=build /go/bin/cmd /app/server

CMD ["/app/server"]
