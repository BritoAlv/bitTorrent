#!/bin/bash

# Check if a container name is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <container_name>"
    exit 1
fi

CONTAINER_NAME=$1

# Check if the container exists
if check_container.sh "$CONTAINER_NAME"; then
    exit 1
fi

# Check if the container is running
if docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    exit 1
else
    echo "Container '${CONTAINER_NAME}' is running."
fi

exit 0