# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Docker parameters
DOCKER=docker
DOCKER_BUILD=$(DOCKER) build
DOCKER_RUN=$(DOCKER) run
DOCKER_IMAGE_NAME=t0mk/rocketreport

# Name of the binary
BINARY_NAME=rocketreport

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

# Run the binary
run:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
	./$(BINARY_NAME)

# Short alias for run
r: run

# Short alias for clean
c: clean

# Short alias for test
t: test

# Build Docker container
docker-build:
	$(DOCKER_BUILD) -t $(DOCKER_IMAGE_NAME) .

# Run Docker container
docker-run:
	$(DOCKER_RUN) -p 8080:8080 $(DOCKER_IMAGE_NAME)

docker-push:
	$(DOCKER) push $(DOCKER_IMAGE_NAME)

.PHONY: all build clean test deps run r c t docker-build docker-run
