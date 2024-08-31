# Variables
DOCKER_IMAGE_NAME := big-john-app
PORT := 5001
ENV_FILE := .env

# Version tagging
VERSION := $(shell git describe --tags --always --dirty)
ifeq ($(VERSION),)
VERSION := dev
endif

# Phony targets
.PHONY: build run run-env stop clean test deps version

# Build the Docker image
build:
	docker build -t $(DOCKER_IMAGE_NAME):$(VERSION) .
	docker tag $(DOCKER_IMAGE_NAME):$(VERSION) $(DOCKER_IMAGE_NAME):latest

# Run the Docker container with individual environment variables
run:
	docker run -p $(PORT):$(PORT) -e LOG_LEVEL=1 -e APP_ENV=development $(DOCKER_IMAGE_NAME):$(VERSION)

# Run the Docker container with environment file
run-env:
	docker run -p $(PORT):$(PORT) --env-file $(ENV_FILE) $(DOCKER_IMAGE_NAME):$(VERSION)

# Stop all running containers of this image
stop:
	docker stop $$(docker ps -q --filter ancestor=$(DOCKER_IMAGE_NAME):$(VERSION))

# Remove the Docker image
clean:
	docker rmi $(DOCKER_IMAGE_NAME):$(VERSION)
	docker rmi $(DOCKER_IMAGE_NAME):latest

postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:16-alpine

# Create a new database named 'bigjohn' inside the running 'bigjohndb' Docker container
# The database is created with 'root' as both the username and owner
createdb:
	docker exec -it postgres16 createdb --username=root --owner=root bigjohn

sqlc:
	sqlc generate

# Run tests (adjust the command as needed for your Go setup)
test:
	go test ./...

# Install or update dependencies
deps:
	go mod tidy

# Build and run in one command (with env file)
up: build run-env

# Stop and remove in one command
down: stop clean

# Print the current version
version:
	@echo $(VERSION)
