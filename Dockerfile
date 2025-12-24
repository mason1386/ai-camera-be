# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies if needed (e.g. for CGO)
# RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/worker cmd/worker/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/migration cmd/migration/main.go

# Production Stage
FROM alpine:latest

WORKDIR /app

# Install basic packages
RUN apk add --no-cache ca-certificates tzdata

# Copy binaries
COPY --from=builder /app/bin/api .
COPY --from=builder /app/bin/worker .
COPY --from=builder /app/bin/migration .

# Copy configurations and assets
COPY config/ config/
COPY migrations/ migrations/
COPY scripts/ scripts/

# Expose API port
EXPOSE 8080

# Default command (can be overridden in docker-compose)
CMD ["./api"]
