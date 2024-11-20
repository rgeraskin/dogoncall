# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY dogoncall/go.mod dogoncall/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY dogoncall/ ./

# Build the Go app
RUN go build .

# Stage 2: Run the Go binary
FROM alpine:3.20

# Create a non-root user and group
RUN addgroup -S app && adduser -S app -G app
USER app

# Set the Current Working Directory inside the container
WORKDIR /app/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/dogoncall .

# Command to run the executable
ENTRYPOINT ["./dogoncall"]