.PHONY: all build test lint clean

all: build

build:
	go build ./...

test:
	go vet ./...
	go test ./...

## lint: run golangci-lint
lint:
	golangci-lint run ./...

clean:
	rm -f qrcode/qrcode qr-gen-server/qr-gen-server
