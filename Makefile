# Variables
DOCKER_IMAGE_NAME := big-john-app
PORT := 5001
ENV_FILE := app.env
DB_URL=postgresql://root:password@localhost:5432/bigjohn?sslmode=disable

# Version tagging
VERSION := $(shell git describe --tags --always --dirty)
ifeq ($(VERSION),)
VERSION := dev
endif

# Phony targets

# Build the Docker image
build:
	docker build -t $(DOCKER_IMAGE_NAME):$(VERSION) .
	docker tag $(DOCKER_IMAGE_NAME):$(VERSION) $(DOCKER_IMAGE_NAME):latest

# Run the Docker container with individual environment variables
run:
	docker run -p $(PORT):$(PORT) --network big-john-network -e APP_ENV=development $(DOCKER_IMAGE_NAME):$(VERSION)

# Run the Docker container with environment file
run-env:
	docker run -p $(PORT):$(PORT) --env-file $(ENV_FILE) $(DOCKER_IMAGE_NAME):$(VERSION)

# Stop all running containers of this image
stop:
	docker stop $$(docker ps -q --filter ancestor=$(DOCKER_IMAGE_NAME):$(VERSION))

# Remove the Docker image
clean:
	docker rmi $(DOCKER_IMAGE_NAME):latest
	docker rmi $(DOCKER_IMAGE_NAME):$(VERSION)

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:16-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root bigjohn

dropdb:
	docker exec -it postgres dropdb bigjohn

migrateup:
	migrate -path internal/db/postgresql/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path internal/db/postgresql/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path internal/db/postgresql/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path internal/db/postgresql/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir internal/db/postgresql/migration -seq $(name)

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

.PHONY: build run run-env postgres version createdb dropdb migrateup migratedown migrateup1 migratedown1 new_migration sqlc 