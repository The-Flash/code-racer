-include .env

NOSOCKET_TARGET = /build/bin/nosocket

dev:
	go run ./cmd/code-racer/main.go

build:
	go build -o bin/code-racer cmd/code-racer/main.go

nosocket:
	CGO_ENABLED=1 go build -o $(NOSOCKET_TARGET) cmd/nosocket/main.go

nosocket-dev:
	go run cmd/nosocket/main.go

docker-build:
	docker build -t $(APP_NAME):latest .
	docker tag $(APP_NAME):latest $(APP_NAME):$(VERSION)

docker-push: docker-build
	docker push $(APP_NAME):latest
	docker push $(APP_NAME):$(VERSION)

compose-build:
	docker-compose build

compose-up:
	docker-compose up

test:
	go test ./...