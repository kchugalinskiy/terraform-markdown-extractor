all: build

build:
	CGO_ENABLED=0 go build -tags netgo -a -o extractor

test-short:
	go test --short ./...