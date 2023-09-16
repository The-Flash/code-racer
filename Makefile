dev:
	go run ./cmd/code-racer/main.go -f manifest.yml -m $(shell pwd)/mntfs

build:
	go build -o code-racer cmd/code-racer/main.go