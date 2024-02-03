# Use the official Go image as the base image
FROM golang:1.16

# Set the working directory inside the container
WORKDIR /app

# Copy the main.go file into the container
COPY main.go .

# Build the Go application
RUN go build -o app

# Expose port 8080 for the HTTP server
EXPOSE 8080

# Set the entry point for the container
ENTRYPOINT ["./app"]
