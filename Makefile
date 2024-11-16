# Makefile

# Var
CONFIG_FILE ?= config/config.yaml
DB_CONTAINER_NAME ?= xmgo_db
DB_IMAGE ?= postgres:13
DB_PORT ?= 25432
DB_USER ?= xmgo
DB_PASSWORD ?= xmgopass
DB_NAME ?= xmgo_db
IMAGE_NAME ?= xmgo
CONTAINER_NAME ?= xmgo_svc
PORT ?= 8080

# Help (default)
.PHONY: help
help: ## Show available commands
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# Run tests
.PHONY: test
test: ## Run tests(only for company_handler_test.go)
	@echo "Running tests..."
	go test -v internal/infrastructure/server/handler/company_handler_test.go

# Run service
.PHONY: run
run: ## Run service
	@echo "Starting service..."
	go run cmd/server/main.go -config $(CONFIG_FILE)

# Start database container
.PHONY: start-db
start-db: ## Start database container
	@echo "Starting database container..."
	docker run --name $(DB_CONTAINER_NAME) -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_DB=$(DB_NAME) -p $(DB_PORT):5432 -d $(DB_IMAGE)

# Stop database container
.PHONY: stop-db
stop-db: ## Stop database containerRE
	@echo "Stopping database container..."
	docker stop $(DB_CONTAINER_NAME)
	docker rm $(DB_CONTAINER_NAME)

# Build docker image
.PHONY: docker-build
docker-build: ## Build docker image
	@echo "Building docker image..."
	docker build -t $(IMAGE_NAME) .

# Run docker container
.PHONY: docker-run
docker-run: ## Run docker container
	@echo "Running docker container..."
	docker run --name $(CONTAINER_NAME) --link $(DB_CONTAINER_NAME):db --link kafka:kafka -e CONFIG_FILE=$(CONFIG_FILE) -e PORT=$(PORT) -p $(PORT):$(PORT) $(IMAGE_NAME)

# Stop docker container
.PHONY: docker-stop
docker-stop: ## Stop docker container
	@echo "Stopping docker container..."
	docker stop $(CONTAINER_NAME)
	docker rm $(CONTAINER_NAME)

# Start kafka container
.PHONY: start-kafka
start-kafka: ## Start kafka container
	@echo "Starting kafka container..."
	docker run -d --name zookeeper -p 2181:2181 zookeeper:3.6.3
	docker run -d --name kafka -p 9092:9092 --link zookeeper:zookeeper \
		-e KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181 \
		-e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 \
		-e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \
		wurstmeister/kafka:2.13-2.7.0
	docker exec -it kafka kafka-topics.sh --create --topic company_events --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1

# Stop kafka container
.PHONY: stop-kafka
stop-kafka: ## Stop kafka container
	@echo "Stopping kafka container..."
	docker stop kafka
	docker stop zookeeper
	docker rm kafka
	docker rm zookeeper

# Restart kafka container
.PHONY: restart-kafka
restart-kafka: stop-kafka start-kafka ## Restart kafka container

# Remove docker images
.PHONY: docker-rmi
docker-rmi: ## Remove docker images
	@echo "Removing docker images..."
	docker rmi xmgo wurstmeister/kafka:2.13-2.7.0 zookeeper:3.6.3