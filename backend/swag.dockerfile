# Dockerfile for Swagger
FROM golang:1.22-alpine

# Install swag CLI tool
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Set the working directory
WORKDIR /app

# Copy the project files
COPY . .

# Generate Swagger documentation
RUN swag init -g backend/main.go

# Expose port 8080 for Swagger UI
EXPOSE 8080

# Start the Swagger UI server
CMD ["swag", "serve", "-p", "8080"]