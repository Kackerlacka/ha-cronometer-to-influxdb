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
RUN go build -o /root/cronapp .

# Stage 2: Create the final image
FROM frolvlad/alpine-glibc:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /root/cronapp /root/cronapp

# Set executable permissions for the cronapp binary
RUN chmod +x /root/cronapp

# Copy the entrypoint script
COPY run.sh /root/

# Set executable permissions for the entrypoint script
RUN chmod +x /root/run.sh

# Run your Go application
CMD ["/root/run.sh"]
