# Build stage
FROM golang:1.23-alpine AS builder

# Install git and ca-certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates tzdata gcc musl-dev sqlite-dev

# Create appuser for security
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the binary with CGO enabled for SQLite support
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s' \
    -a -installsuffix cgo \
    -o cloudgate-backend .

# Final stage - use alpine instead of scratch for CGO compatibility
FROM alpine:latest

# Install sqlite and ca-certificates
RUN apk add --no-cache ca-certificates sqlite tzdata

# Import from builder
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary
COPY --from=builder /build/cloudgate-backend /cloudgate-backend

# Use non-root user
USER appuser

# Expose the port defined by environment variable (Cloud Run uses 8080 by default)
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/cloudgate-backend"] 