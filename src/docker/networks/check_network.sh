#!/bin/bash

# Check if a Docker network exists
NETWORK_NAME=$1

if [ -z "$NETWORK_NAME" ]; then
    echo "Usage: $0 <network_name>"
    exit 1
fi

if docker network ls --format '{{.Name}}' | grep -wq "$NETWORK_NAME"; then
    echo "Network '$NETWORK_NAME' exists."
    exit 0
else
    echo "Network '$NETWORK_NAME' does not exist."
    exit 1
fi