#!/bin/bash

# Check if a container name is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <container_name>"
    exit 1
fi

CONTAINER_NAME=$1

# Check if a container exists
if [ "$(docker ps -a -q -f name=^/${CONTAINER_NAME}$)" ]; then
    echo "Container '$CONTAINER_NAME' exists."
    exit 0
else
    echo "Container '$CONTAINER_NAME' does not exist."
    exit 1
fi