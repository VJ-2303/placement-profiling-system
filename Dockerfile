# KCT Placement Profiling System - Backend
# Build from repository root, targeting backend folder

FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git (required for fetching dependencies)
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files from backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source code
COPY backend/ .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o main ./cmd/api

# Final stage - minimal image
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy migrations folder
COPY --from=builder /app/migrations ./migrations

# Create non-root user for security
RUN adduser -D -g '' appuser
USER appuser

# Expose port
EXPOSE 4000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:4000/health || exit 1

# Run the binary
CMD ["./main"]
