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

# Create the final image with MySQL and the Go application
FROM mysql:8.0

# Set MySQL environment variables
ENV MYSQL_ROOT_PASSWORD=password4321
ENV MYSQL_DATABASE=blueprint
ENV MYSQL_USER=melkey
ENV MYSQL_PASSWORD=password1234


# Initialize the MySQL data directory
RUN ["mysqld", "--initialize-insecure", "--datadir=/var/lib/mysql"]

# Expose MySQL and app ports
EXPOSE 3306 8080

# Copy the built Go application from the builder stage
COPY --from=builder /app/main /main
# Add a custom SQL file for setting permissions and initializing data
COPY init.sql /docker-entrypoint-initdb.d/

# Command to start MySQL and your Go application
CMD ["sh", "-c", "mysqld --bind-address=0.0.0.0 & sleep 10 && ./main"]
