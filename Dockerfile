# Use the official Golang image as the base image
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application code to the working directory
COPY . .

# Ensure all dependencies are included
RUN go mod tidy

# Build the Go application
RUN go build -o xmgo ./cmd/server

# Use a minimal base image for the final stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go application from the builder stage
COPY --from=builder /app/xmgo .
COPY --from=builder /app/config/config.yaml ./config/config.yaml
COPY --from=builder /app/migrations ./migrations

# Set environment variables for the config file and port
ENV CONFIG_FILE=/app/config/config.yaml
ENV PORT=8080

# Expose the application port using an environment variable
EXPOSE ${PORT}

# Command to run the application with configurable port and config file
CMD ["sh", "-c", "./xmgo -config=$CONFIG_FILE"]