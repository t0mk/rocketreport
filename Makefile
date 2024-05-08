VERSION_FILE=version
VERSION := $(shell cat ${VERSION_FILE})

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

NOW=${shell date +%F_%T}

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

plugins_md:
	$(GOCMD) run cmd/generate_PLUGINS.md/main.go

# Build the binary for local testing
build: plugins_md
	$(GOBUILD) -ldflags "-X main.buildTime=$(NOW) -X main.version=$(VERSION)" -o $(BINARY_NAME) $(MAIN_PATH)

static-build-amd64:
	# build statically compiled binary
	# there's ld warning "/usr/bin/ld: warning: bint-x64-amd64.o: missing .note.GNU-stack section implies executable stack" and I don't know how to get rid of it
	${DOCKER_RUN} --rm -v ${shell pwd}:/app ${DOCKER_BUILDER_IMAGE_NAME}:${VERSION} /app/build-inside-container.sh amd64 || true

static-build-arm64:
	# build statically compiled binary
	${DOCKER_RUN} --rm -v ${shell pwd}:/app ${DOCKER_BUILDER_IMAGE_NAME}:${VERSION} /app/build-inside-container.sh arm64

static-builds: static-build-amd64 static-build-arm64

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

docker-builder-image:
	# build builder image
	${DOCKER_BUILD} -t ${DOCKER_BUILDER_IMAGE_NAME}:${VERSION} -f docker/builder .

docker-image: static-build-amd64
	@if ! test `find ${BINARY_NAME_AMD64} -newermt "10 seconds ago"`; then echo "binary was not created in last 5 secs"; exit 1; fi
	# build actual docker image for t0mk/rocketreport
	$(DOCKER_BUILDX) build --platform=linux/amd64 -t $(DOCKER_IMAGE_NAME) -f docker/Dockerfile --load .


docker-check:
	${DOCKER_RUN} --rm ${DOCKER_IMAGE_NAME} version
	${DOCKER_RUN} --rm ${DOCKER_IMAGE_NAME} plugin gasPrice

docker-push:
	$(DOCKER) push $(DOCKER_IMAGE_NAME)

.PHONY: all build static-build clean test deps c t docker-builder-image docker-image docker-check docker-push
