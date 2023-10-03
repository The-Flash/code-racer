NOSOCKET_INSTALL_PATH = /usr/local/bin
NOSOCKET_TARGET = tmp/nosocket

dev:
	go run ./cmd/code-racer/main.go -f manifest.yml -m $(shell pwd)/mntfs -r $(shell pwd)/runners

build:
	go build -o bin/code-racer cmd/code-racer/main.go

nosocket:
	go build -o $(NOSOCKET_TARGET) cmd/nosocket/main.go

nosocket-dev:
	go run cmd/nosocket/main.go

install_nosocket:
	make nosocket
	mv $(NOSOCKET_TARGET) $(NOSOCKET_INSTALL_PATH)