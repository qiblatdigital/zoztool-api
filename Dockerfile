# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/api ./cmd/api

# Final stage
FROM alpine:3.20

WORKDIR /app

# Install ca-certs and curl (for healthcheck)
RUN apk add --no-cache ca-certificates curl

# Copy binary
COPY --from=builder /app/bin/api .

EXPOSE 8083

CMD ["./api"]
