build:
	go build ./...

test:
	go test ./...

fmt:
	go fmt ./...

.PHONY: build test fmt
