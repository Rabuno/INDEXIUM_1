# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build binary từ file main.go nằm trong cmd/server
RUN go build -o main cmd/server/main.go

# Stage 2: Run
FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /app/main .
# Cài thêm curl nếu cần debug healthcheck (tùy chọn)
RUN apk --no-cache add curl
CMD ["./main"]