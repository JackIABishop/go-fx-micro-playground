

# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy go mod files and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source code (to resolve internal/logging, etc.)
COPY . .

# Set working directory to the gateway service
WORKDIR /app/services/gateway
RUN go build -o gateway .

# Final runtime stage
FROM alpine:latest
WORKDIR /root/

# Copy the compiled binary
COPY --from=builder /app/services/gateway/gateway .

# Expose the service port
EXPOSE 8080

# Healthcheck for container orchestration
HEALTHCHECK --interval=10s --timeout=5s --retries=3 \
  CMD wget --no-verbose --spider http://localhost:8080/health || exit 1

# Run the service
CMD ["./gateway"]