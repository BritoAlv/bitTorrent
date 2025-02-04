#!/bin/bash

# Check if the correct number of arguments is provided
if [ "$#" -ne 2 ]; then
    echo "❌ Usage: $0 <container_name_or_id> <network_name>"
    exit 1
fi

CONTAINER_NAME_OR_ID=$1
# Check if the container exists
if ! docker ps -a --format '{{.Names}}' | grep -wq "$1"; then
    echo "❌ Container $1 does not exist."
    exit 1
fi

# Check if the network exists
if ! docker network ls --format '{{.Name}}' | grep -wq "$2"; then
    echo "❌ Network $2 does not exist."
    exit 1
fi

CONTAINER_NAME_OR_ID=$1
NETWORK_NAME=$2

echo "🔍 Checking if the container is connected to the network..."

# Disconnect the container from the network
if docker network disconnect "$NETWORK_NAME" "$CONTAINER_NAME_OR_ID"; then
    echo "✅ Successfully disconnected $CONTAINER_NAME_OR_ID from $NETWORK_NAME"
else
    echo "❌ Failed to disconnect $CONTAINER_NAME_OR_ID from $NETWORK_NAME"
    exit 1
fi