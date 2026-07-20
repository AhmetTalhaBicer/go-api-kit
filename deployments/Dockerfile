# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Pre-copy dependencies for layer caching
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Compile binary — CGO disabled for static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/build/server ./cmd/api/...

# ---

# Runtime stage — minimal alpine image
FROM alpine:3.20

# Security: create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy binary only
COPY --from=builder /app/build/server .

USER appuser

EXPOSE 8080

ENTRYPOINT ["./server"]
