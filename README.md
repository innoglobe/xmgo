
# XMGO Service

This is the XMGO service, which provides functionality for managing companies and sending events to Kafka.

## Requirements

- Go 1.23 or higher
- Docker

## Configuration

The service configuration is stored in the `config/config.yaml` file. Ensure that the Kafka broker address and database settings are correctly configured.

## Running the Service(from docker)

1. **Start the database container**:
   ```make start-db```
2. **Start the Kafka container**:
   ```make start-kafka```
3. **Build the Docker image**:  
   ```make docker-build```
4. **Run the Docker container**:  
   ```make docker-run```

## Running Locally(using go run)
1. **Start the database container**:
   ```make start-db```
2. **Start the Kafka container**:
   ```make start-kafka```
3. **Run the service**:  
   ```make run``` OR ```go run cmd/server/main.go -config config/config.yaml```

## Integration Tests
You run the integration tests(specifically for the company handler) using:
```make test``` or ```go test -v internal/infrastructure/server/handler/company_handler_test.go```
It is recommended to run ```make stop-db``` and ```make start-db``` in order to reset the database each time you want to run the tests.

## Makefile Targets
- ```help```: Show available commands
- ```test```: Run tests (only for company_handler_test.go)
- ```run```: Run the service locally
- ```start-db```: Start the database container
- ```stop-db```: Stop the database container(with rm)
- ```docker-build```: Build the Docker image
- ```docker-run```: Run the Docker container
- ```docker-stop```: Stop the Docker container
- ```start-kafka```: Start the Kafka container
- ```stop-kafka```: Stop the Kafka container
- ```restart-kafka```: Restart the Kafka container
- ```docker-rmi```: Remove Docker images(removes xmgo, zookeeper and kafka images)

## Configuration File
The ```config/config.yaml``` file contains the following settings:

    production: false
    server:
      host: 0.0.0.0
      port: 8080
      ssl:
        enabled: false
        cert_file: ""
        key_file: ""
      timeout: 5
    
    database:
      # for docker-run host should be db and port 5432
      # for go run host should be 127.0.0.1 and port 25432
      host: 127.0.0.1
      port: 25432
      user: xmgo
      pass: xmgopass
      name: xmgo_db
    
    jwt:
      secret: secret-key
    
    kafka:
      brokers:
        - "localhost:9092"
      topic: "company_events"

Ensure that the ```host``` and ```port``` settings for the database and Kafka are correctly configured based on whether you are running the service locally or in Docker.