top: build gev

gev:
	go build -o gev ./cmd/eval

build: codegen
	go build ./...

codegen:
	make -C 
test:
	go test ./...

fmt:
	go fmt ./...

.PHONY: build test fmt
