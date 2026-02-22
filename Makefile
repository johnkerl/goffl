build: codegen
	go build ./...

codegen:
	make -C cmd/eval/generated

test:
	go test ./...

fmt:
	go fmt ./...

.PHONY: build codegen test fmt
