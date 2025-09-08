# Build stage
FROM golang:1.24-alpine AS builder

# Install ca-certificates and build tools
RUN apk --no-cache add ca-certificates git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Install Goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy source code
COPY . .

# Build the application
RUN go build -o social-api cmd/api/main.go

# Development stage
FROM golang:1.24-alpine AS development

# Install ca-certificates and build tools
RUN apk --no-cache add ca-certificates git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Install Goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy source code
COPY . .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/social-api .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Copy Goose binary
COPY --from=builder /go/bin/goose ./goose

# Expose port
EXPOSE 8080

# Command to run the application
CMD ["./social-api"]