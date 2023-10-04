FROM golang:1.21.1-alpine3.17
RUN apk update && \
    apk add make
RUN go install github.com/cosmtrek/air@latest

WORKDIR /code-racer
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .