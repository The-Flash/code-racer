FROM golang:1.19-alpine AS build
WORKDIR /build

COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/
COPY Makefile Makefile
COPY go.mod go.mod
COPY go.sum go.sum

RUN apk add --no-cache make

RUN make build

FROM alpine:latest as final
COPY --from=build /build/code-racer /bin/code-racer
COPY manifest.yml manifest.yml
COPY runners/ runners/