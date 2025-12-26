# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install required packages
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files from backend directory
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy all backend source code
COPY backend/cmd ./cmd
COPY backend/internal ./internal
COPY backend/migrations ./migrations

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main ./cmd/api

# Final stage - minimal image
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS and tzdata for timezones
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 4000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:4000/health || exit 1

# Run the application
CMD ["./main"]