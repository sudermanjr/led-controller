# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=led-controller
COMMIT := $(shell git rev-parse HEAD)
VERSION := "dev"

build:
	docker run --rm -ti -v "$(pwd)":/go/src/github.com/sudermanjr/led-controller rpi-ws281x-go-builder /usr/bin/qemu-arm-static /bin/sh -c "go build -o src/github.com/sudermanjr/led-controller -v led-controller"

test:
	printf "Linter:\n"
	GO111MODULE=on $(GOCMD) list ./... | xargs -L1 golint | tee golint-report.out
	printf "\n\nTests:\n\n"
	GO111MODULE=on $(GOCMD) test -v --bench --benchmem -coverprofile coverage.txt -covermode=atomic ./...
	GO111MODULE=on $(GOCMD) vet ./... 2> govet-report.out
	GO111MODULE=on $(GOCMD) tool cover -html=coverage.txt -o cover-report.html
	printf "\nCoverage report available at cover-report.html\n\n"
tidy:
	$(GOCMD) mod tidy
clean:
	$(GOCLEAN)
	$(GOCMD) fmt ./...
	rm -f $(BINARY_NAME)
