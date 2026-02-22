# syntax=docker/dockerfile:1

# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /src

# Install runtime dependencies needed for module fetching (git) and certificates
RUN apk add --no-cache ca-certificates git

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy only the source tree needed to build.
# Adjust these COPY lines to match your repository layout.
# Common top-level folders: cmd/, internal/, pkg/, api/, configs/, etc.
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY config/ ./config/
COPY infrastructure/ ./infrastructure/

# If you have other top-level directories, add explicit COPY for them.
# Avoid `COPY . .` to prevent accidentally adding secrets or dev files.

# Build a static Linux binary
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Change the path to your main package if different
RUN go build -o main cmd/server/main.go

# Stage 2: Runtime
FROM alpine:3.19
# Create a non-root user
RUN addgroup -S app && adduser -S -G app app
WORKDIR /app

# Copy binary from builder
COPY --from=builder /src/main /usr/local/bin/main

# Set ownership and make executable
RUN chown app:app /usr/local/bin/main && chmod +x /usr/local/bin/main
USER app

# Expose port (adjust if your app uses a different port)
EXPOSE 8080

# Optional: if you need curl for debugging/healthcheck, uncomment the next line
# RUN apk add --no-cache curl

ENTRYPOINT ["/usr/local/bin/main"]