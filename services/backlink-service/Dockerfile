# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy shared module first (for replace directive)
COPY shared/go ./shared/go

# Copy service files
WORKDIR /build/services/backlink-service
COPY services/backlink-service/go.mod ./
COPY services/backlink-service/go.su[m] ./
RUN go mod download

# Copy source code
COPY services/backlink-service/ .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/backlink-service \
    ./cmd/main.go

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -s /bin/sh -D appuser

# Copy binary from builder
COPY --from=builder /app/backlink-service .

# Copy migrations
COPY --from=builder /build/services/backlink-service/migrations ./migrations

# Set ownership
RUN chown -R appuser:appgroup /app

USER appuser

# Expose port
EXPOSE 8082

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8082/health || exit 1

# Run
CMD ["./backlink-service"]
