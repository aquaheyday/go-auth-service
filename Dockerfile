# Build Stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -o auth-server \
    ./cmd/auth-server

# Runtime Stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /app/auth-server .
COPY configs/config.yaml .
EXPOSE 50051
CMD ["./auth-server"]