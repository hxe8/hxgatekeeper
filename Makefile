APP=gatekeeper

.PHONY: run build test fmt

run:
	go run ./cmd/$(APP)

build:
	go build -o bin/$(APP) ./cmd/$(APP)

test:
	go test ./...

fmt:
	gofmt -w $(shell find . -name '*.go' -type f)
