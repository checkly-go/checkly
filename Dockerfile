# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the server application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/server .

# Set environment variable for MongoDB URI (can be overridden at runtime)
ENV MONGO_URI="mongodb+srv://hawkaii:hawkaii2022@cluster0.cqikohy.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

# Run the server application
CMD ["./server"]