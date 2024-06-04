.PHONY: build generate

build:
	go build -o bin/app cmd/app/*.go

generate:
	go generate ./...
