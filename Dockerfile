FROM golang:1.19-alpine AS build
RUN apk add --no-cache make

WORKDIR /build
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/
COPY Makefile Makefile
COPY go.mod go.mod
COPY go.sum go.sum


RUN make build

FROM alpine:latest as final
COPY --from=build /build/code-racer /bin/code-racer
COPY runners/ runners/