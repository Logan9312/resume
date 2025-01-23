# Use the official Go 1.23 Alpine image
FROM golang:1.23-alpine

# Install system dependencies, including poppler-utils
RUN apk add --no-cache poppler-utils

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files first for dependency caching
COPY go.mod ./
COPY go.sum ./

# Download and cache Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . ./

# Build the application
RUN go build -o main .

# Expose the port your app will run on
EXPOSE 3000

# Command to run the application
CMD ["./main"]
