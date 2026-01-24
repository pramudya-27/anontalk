# Start from the official Golang image to build the binary
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=0 creates a statically linked binary (no libc dependency), small and safe for scratch/alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Start a new stage from a small Alpine image
FROM alpine:latest

# Install ca-certificates (needed for HTTPS)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .
# Copy migrations if needed, though usually better handled by a separate init container or tool
COPY --from=builder /app/migrations ./migrations
# Copy uploads folder structure locally (optional, but good for structure)
RUN mkdir -p uploads

# Expose the API port
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./main"]
# Default arguments (can be overridden in K8s)
CMD ["-mode", "all"]
