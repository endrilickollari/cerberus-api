# Build stage
FROM golang:1.23.0 AS builder

LABEL authors="_endrilickollari"

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /cerberus ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /cerberus /app/cerberus

# Expose port and set environment variables
ENV PORT=8080
EXPOSE $PORT

# Run the app
CMD ["/app/cerberus"]

