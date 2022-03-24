SHELL := /bin/bash

CONTAINER_RUNTIME ?= podman
IMAGE_TAG := 0.0.1
IMAGE_NAME := anjannath/raisepr

GO_SOURCE := pkg
all: raisepr

raisepr: main.go $(GO_SOURCE)
	go build -o $@ $<

.PHONY: clean
clean:
	rm -rf raisepr

.PHONY: test
test:
	@go test -v ./...

container-build: Dockerfile
	$(CONTAINER_RUNTIME) build -f Dockerfile -t $(IMAGE_NAME):$(IMAGE_TAG) .
