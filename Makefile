VERSION_FILE=version
VERSION := $(shell cat ${VERSION_FILE})

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Docker parameters
DOCKER=docker
DOCKER_BUILD=$(DOCKER) build
DOCKER_BUILDX=$(DOCKER) buildx
DOCKER_RUN=$(DOCKER) run
DOCKER_IMAGE_NAME=t0mk/rocketreport
DOCKER_BUILDER_IMAGE_NAME=t0mk/rocketreport

# Name of the binary
BINARY_NAME=rocketreport
BINARY_NAME_AMD64=${BINARY_NAME}-amd64

# Main package path
MAIN_PATH=./cmd/rocketreport/main.go

# Default target
all: test build

docs:
	$(GOCMD) run cmd/generate_readme/main.go

# Build the binary
build: docs
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

# Clean the binary
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Test the project
test:
	$(GOTEST) -v ./...

# Install dependencies
deps:
	$(GOGET) ./...

# Short alias for clean
c: clean

# Short alias for test
t: test

# Build Docker container
docker-build:
	# build builder image
	${DOCKER_BUILD} -t ${DOCKER_BUILDER_IMAGE_NAME}:${VERSION} -f docker/builder .
	# build statically compiled binary
	${DOCKER_RUN} --rm -v ${shell pwd}:/app ${DOCKER_BUILDER_IMAGE_NAME}:${VERSION} /app/build-inside-container.sh
	@if ! test `find ${BINARY_NAME_AMD64} -newermt "10 seconds ago"`; then echo "binary was not created in last 5 secs"; exit 1; fi
	# build actual docker image for t0mk/rocketreport
	$(DOCKER_BUILDX) build --platform=linux/amd64 -t $(DOCKER_IMAGE_NAME) --load .

docker-testrun:
	${DOCKER_RUN} --rm ${DOCKER_IMAGE_NAME} plugin gasPrice

docker-push:
	$(DOCKER) push $(DOCKER_IMAGE_NAME)

.PHONY: all build clean test deps run r c t docker-build docker-run docker-push
