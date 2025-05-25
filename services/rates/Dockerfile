

# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Cache module downloads
COPY go.mod ./
RUN go mod download

# Copy service code and build
COPY . .
WORKDIR /app/services/rates
RUN go build -o rates .

# Final runtime stage
FROM alpine:latest
WORKDIR /root/

# Copy the compiled binary
COPY --from=builder /app/services/rates/rates .

# Expose the service port
EXPOSE 8081

# Healthcheck for container orchestration
HEALTHCHECK --interval=10s --timeout=5s --retries=3 \
  CMD wget --no-verbose --spider http://localhost:8081/health || exit 1

# Run the service
CMD ["./rates"]