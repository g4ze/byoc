#!/bin/sh

echo "parent container is up and running"

# Check if required environment variables are set
if [ -z "$DOCKER_IMAGE_NAME" ] ||[ -z "$PORT"]; then
  echo "DOCKER_IMAGE_NAME and PORT environment variables must be set."
  exit 1
else  
  echo "DOCKER_IMAGE_NAME: $DOCKER_IMAGE_NAME PORT: $PORT"
fi
sudo service docker start
docker pull "$DOCKER_IMAGE_NAME" # from docker hub or wherever the image is stored
  # Confirm the image is built
    if [ $? -eq 0 ]; then
    echo "Docker image $DOCKER_IMAGE_NAME pulled successfully."
    else
    echo "Failed to pull Docker image $DOCKER_IMAGE_NAME."
    exit 1
    fi
    docker run -d -p "$PORT:$PORT" "$DOCKER_IMAGE_NAME"
    if [ $? -eq 0 ]; then
    echo "Container is up and running, accessible at port $PORT."
    else
    echo "Failed to start the Docker container."
    exit 1
  fi
tail -f /dev/null






