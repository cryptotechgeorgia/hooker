# Project variables
APP_NAME := thooker
DOCKER_COMPOSE_FILE := docker-compose.build.yml
DOCKER_REGISTRY := ###

# Targets
.PHONY: all build run test clean docker-build docker-up docker-down format vet lint

build:
	@echo "Building the Go project..."
	@go build -o $(APP_NAME) cmd/main.go

run: build
	@echo "Running the Go project..."
	@./$(APP_NAME)

deploy:
	@echo "Building Docker image..."
	@read -p "Enter Docker image tag : " IMAGE_TAG; \
	docker build -t $(DOCKER_REGISTRY)/$(APP_NAME):$$IMAGE_TAG . ;\
       	docker push  $(DOCKER_REGISTRY)/$(APP_NAME):$$IMAGE_TAG
