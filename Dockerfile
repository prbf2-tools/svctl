# Build stage
FROM golang:1.24.1-alpine AS builder

# Install required tools for building
RUN apk add --no-cache git make protoc protobuf-dev

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY internal internal
COPY pkg pkg
COPY cmd cmd
COPY svctl svctl
COPY main.go main.go

COPY Makefile ./

# Generate code and build the daemon
RUN make build

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests and Docker CLI for Docker game servers
RUN apk add ca-certificates docker-cli

# Copy binary from build stage
COPY --from=builder /build/bin/svctl /usr/local/bin/svctl

# Set working directory
WORKDIR /app

# Expose default daemon port
EXPOSE 9090

# Default command to run daemon
CMD ["svctl", "daemon", "--data-dir", "/app/data", "--config-dir", "/app/config"]
