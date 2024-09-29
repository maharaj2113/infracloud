# Use an official Golang runtime as a parent image
FROM golang:1.21-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go application
RUN go build -o /url-shortener

# Make port 8080 available to the world outside this container
EXPOSE 8080

# Run the Go app
CMD ["/url-shortener"]