# Use an official Go image as the base image
FROM golang:1.23 as builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose MySQL and app ports
EXPOSE 8080

# Copy the built Go application from the builder stage
COPY --from=builder /app/main /main

# Command to start MySQL and your Go application
CMD ["/main"]
