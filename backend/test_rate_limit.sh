#!/bin/bash

echo "Testing rate limiter with rapid requests..."
echo "Making 5 rapid requests to test rate limiting (should see 429 error after 3 requests)..."
echo

for i in {1..5}
do
    echo "Request $i:"
    curl -s -w "\nHTTP Status: %{http_code}\n" \
        -H "Content-Type: application/json" \
        -d '{"username": "testuser", "password": "Password123!"}' \
        http://localhost:8000/auth/login
    echo "-------------------"
    sleep 0.1  # Reduced delay to test rate limiting
done
