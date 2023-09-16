FROM golang:1.19-alpine AS build
WORKDIR /build

COPY cmd/ cmd/
COPY internal/ internal/
COPY Makefile Makefile
COPY go.mod go.mod
COPY go.sum go.sum

RUN make build

FROM alpine:latest as final
COPY --from=build /build/code-racer /bin/code-racer
COPY manifest.yml manifest.yml
ENTRYPOINT [ "/bin/code-racer", "-f", "manifest.yml", "-m", "/opt/code-racer/" ]