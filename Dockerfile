FROM golang:1.21 AS build
# install system dependencies
RUN apt-get update
RUN apt-get install -y build-essential
RUN apt-get install libseccomp-dev
RUN apt-get install make
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
# build nosocket binary to /build/bin/nosocket
RUN make install_nosocket
COPY ./entrypoint.sh ./entrypoint.sh

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
RUN make install_nosocket

FROM debian:12 as final
COPY --from=build /build/bin/code-racer /bin/code-racer
COPY --from=build /build/bin/nosocket /bin/nosocket
COPY --from=build /build/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
COPY runners/ runners/