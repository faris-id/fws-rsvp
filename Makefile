.PHONY: all test mod tidy run

all: run

test:
	go test ./...

mod:
	go mod download
	go mod verify

tidy:
	go mod tidy

run :
	go run app/web-service/main.go
