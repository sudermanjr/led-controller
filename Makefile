# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=led-controller
COMMIT := $(shell git rev-parse HEAD)
VERSION=$(shell git describe --tags)
HYPRIOT_IMAGE="https://github.com/hypriot/image-builder-rpi/releases/download/v1.12.0/hypriotos-rpi-v1.12.3.img.zip"
PKG_PATH=/go/src/github.com/sudermanjr/led-controller
DOCKER_GOCACHE=/root/.cache/go-build
BUILD_IMG=ws2811-builder:latest
PLATFORM=linux/arm/v6
LOCAL_TMP=$(PWD)/.tmp

all: lint test create-builder build
build: create-builder build-local-arm
build-local-arm: create-builder
	docker run --rm -v "${PWD}":$(PKG_PATH) --platform $(PLATFORM) -w "$(PKG_PATH)" $(BUILD_IMG) $(GOBUILD) -ldflags="-X main.version=$(VERSION) -X main.commit=$(COMMIT) -s -w"
build-osx:
	$(GOBUILD)
create-builder:
	docker buildx build --platform $(PLATFORM) --tag $(BUILD_IMG) .
lint:
	golangci-lint run
reportcard:
	goreportcard-cli -t 100 -v
test:
	$(GOCMD) test -v --bench --benchmem -coverprofile coverage.txt -covermode=atomic ./...
	$(GOCMD) vet ./... 2> govet-report.out
	$(GOCMD) tool cover -html=coverage.txt -o cover-report.html
	printf "\nCoverage report available at cover-report.html\n\n"
tidy:
	$(GOCMD) mod tidy
clean:
	$(GOCLEAN)
	$(GOCMD) fmt ./...
	rm -f $(BINARY_NAME)
	rm -f pi-builder/init-decrypted
flash:
	sops -d pi-builder/cloud-init.yaml | yq r - data > pi-builder/init-decrypted
	flash --bootconf pi-builder/no-uart-config.txt --userdata pi-builder/init-decrypted $(HYPRIOT_IMAGE)
	rm pi-builder/init-decrypted
encrypt-init:
	sops --encrypt pi-builder/init-decrypted > pi-builder/cloud-init.yaml
decrypt-init:
	sops -d pi-builder/cloud-init.yaml | yq r - data > pi-builder/init-decrypted
styleguide:
	stylemark -i pkg/dashboard/assets/css -o stylemark -c .stylemark.yml -w -b 8081
run-dashboard:
	go run main.go dashboard
