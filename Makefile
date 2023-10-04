NOSOCKET_TARGET = /build/bin/nosocket

dev:
	go run ./cmd/code-racer/main.go -f manifest.yml -m $(shell pwd)/mntfs -r $(shell pwd)/runners

build:
	go build -o bin/code-racer cmd/code-racer/main.go

nosocket:
	CGO_ENABLED=1 go build -o $(NOSOCKET_TARGET) cmd/nosocket/main.go

nosocket-dev:
	go run cmd/nosocket/main.go