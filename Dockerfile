# Start from the official Go image
FROM golang:1.23-alpine AS builder

# Install git for fetching dependencies
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest  

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy the uploads folder structure if needed (optional, depends on your app)
# RUN mkdir -p /uploads/t_hub_document/employee

# Expose port 9005 to the outside world
EXPOSE 9005

# Command to run the executable
CMD ["./main"]
