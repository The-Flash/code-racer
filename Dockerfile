FROM golang:1.21.1-alpine3.17 AS build
# install system dependencies
RUN apk add make
WORKDIR /build
# copy project files
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/
COPY Makefile Makefile
COPY go.mod go.mod
COPY go.sum go.sum
# build binary to /build/bin/code-racer
RUN make build
COPY ./entrypoint.sh ./entrypoint.sh

FROM --platform=$BUILDPLATFORM golang:1.21.1-alpine3.17 as nosocketbuild
ARG TARGETOS TARGETARCH
RUN apk add clang lld
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
ENV CGO_ENABLED=1

RUN apk add pkgconfig musl-dev gcc libseccomp-dev make

WORKDIR /build
COPY cmd/ cmd/
COPY Makefile Makefile
COPY go.mod go.mod
COPY go.sum go.sum
# build nosocket binary to /build/bin/nosocket
RUN make nosocket

FROM golang:1.21 AS dev
RUN apt-get update
RUN apt-get install -y build-essential
RUN apt-get install libseccomp-dev
RUN apt-get install make
RUN go install github.com/cosmtrek/air@latest
WORKDIR /code-racer
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN make nosocket

FROM debian:12 as final
COPY --from=build /build/bin/code-racer /bin/code-racer
COPY --from=nosocketbuild /build/bin/nosocket /bin/nosocket
COPY --from=build /build/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
COPY runners/ runners/