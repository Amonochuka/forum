#!/bin/bash

set -e

IMAGE_NAME="go-forum"
CONTAINER_NAME="go-forum-app"

echo "Building Docker image..."
docker build -t $IMAGE_NAME .

echo "Removing old container (if exists)..."
docker rm -f $CONTAINER_NAME 2>/dev/null || true

echo "Running container..."

docker run -d \
  --name $CONTAINER_NAME \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  $IMAGE_NAME

echo "App running at http://localhost:8080"