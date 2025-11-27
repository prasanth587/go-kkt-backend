# Use the official Golang image as the base image
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files (go.mod and go.sum) to the working directory
COPY go.mod go.sum ./

# Download Go modules (dependencies)
RUN go mod download

# Copy the application source code to the working directory
COPY . .

# Build the Go application
RUN go build -o go-transport-hub .

# Create a new lightweight image for running the application
FROM golang:1.23-alpine

# Set the working directory
WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/go-transport-hub .

# Expose the port that the application listens on (if applicable)
EXPOSE 9005

# Command to run the executable
CMD ["./go-transport-hub"]
