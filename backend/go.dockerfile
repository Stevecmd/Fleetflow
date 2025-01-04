FROM golang:1.22-alpine

# Set the working directory
WORKDIR /app

# Install git for downloading dependencies
RUN apk add --no-cache git

# Copy env file
COPY .env ./

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy the rest of the code
COPY . .

RUN go build -o main

# Expose the port
EXPOSE 8000

# Set the entry point
CMD ["./main"]
