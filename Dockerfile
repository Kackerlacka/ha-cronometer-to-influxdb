# Stage 1: Build the Go executable
FROM golang:latest AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o go-app .

# Stage 2: Create the final image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/go-app .

# Install cron
RUN apk add --no-cache dcron

# Copy the entrypoint script
COPY run.sh .

# Set executable permissions on run.sh
RUN chmod +x run.sh

# Command to run the entrypoint script
CMD ["./run.sh"]
