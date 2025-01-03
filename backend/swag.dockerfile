# Dockerfile for Swagger
FROM swaggerapi/swagger-ui:latest

# Copy the swagger.json file
COPY backend/docs/swagger.json /usr/share/nginx/html/swagger.json

# Set environment variables for Swagger UI
ENV SWAGGER_JSON=/usr/share/nginx/html/swagger.json
ENV BASE_URL=/swagger

# Expose port 8080
EXPOSE 8080