# Variables
DOCKER_IMAGE_NAME := big-john-app
DOCKER_TAG := latest
PORT := 5001
ENV_FILE := .env

# Phony targets
.PHONY: build run run-env stop clean test deps

# Build the Docker image
build:
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) .

# Run the Docker container with individual environment variables
run:
	docker run -p $(PORT):$(PORT) -e LOG_LEVEL=1 -e APP_ENV=development $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

# Run the Docker container with environment file
run-env:
	docker run -p $(PORT):$(PORT) --env-file $(ENV_FILE) $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

# Stop all running containers of this image
stop:
	docker stop $$(docker ps -q --filter ancestor=$(DOCKER_IMAGE_NAME):$(DOCKER_TAG))

# Remove the Docker image
clean:
	docker rmi $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

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
