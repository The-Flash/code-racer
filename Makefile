dev:
	go run ./cmd/code-racer/main.go -f manifest.yml -m $(shell pwd)/mntfs -r $(shell pwd)/runners

build:
	go build -o code-racer cmd/code-racer/main.go