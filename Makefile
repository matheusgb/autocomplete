.PHONY: run test lint build up down

run:
	go run main.go

test:
	go test -v ./...

lint:
	golangci-lint run ./...

build:
	go build -o autocomplete .

up:
	docker compose up --build

down:
	docker compose down
