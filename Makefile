build:
	@go build -o bin/exporter

run: build
	@./bin/exporter

test:
	@go test -v ./...