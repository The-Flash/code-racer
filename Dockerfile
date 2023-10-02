FROM golang:1.21 AS build
RUN apk add --no-cache gcc
RUN apk add --no-cache libseccomp-dev
RUN apk add --no-cache make
WORKDIR /build
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/
COPY Makefile Makefile
COPY go.mod go.mod
COPY go.sum go.sum
RUN make build

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

FROM alpine:latest as final
COPY --from=build /build/code-racer /bin/code-racer
COPY runners/ runners/