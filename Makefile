.PHONY: build docker

TARGET_DIR ?= ./

# TEST_ARGS = -v | grep -c RUN
DOCKER_OPTS ?= --rm

VERSION := $(shell git describe --tags --abbrev=0)

TAG ?= $(VERSION)

build:
	go build -o $(TARGET_DIR) ./cmd/user-stats-service

docker:
	docker build -t github.com/grippenet/user-stats-service:$(VERSION)  -f build/docker/Dockerfile $(DOCKER_OPTS) .