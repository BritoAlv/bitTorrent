#!/bin/bash

# Function to print usage
usage() {
    echo "Usage: $0 <container_name> <network_name>"
    exit 1
}

# Check if the correct number of arguments is provided
if [ "$#" -ne 2 ]; then
    usage
fi

CONTAINER_NAME=$1
NETWORK_NAME=$2

# Check if the container exists and is running
echo "🔍 Checking if the container '$CONTAINER_NAME' exists and is running..."
if ! docker ps --format '{{.Names}}' | grep -w "$CONTAINER_NAME"; then
    echo "❌ Container '$CONTAINER_NAME' does not exist or is not running."
    exit 1
fi
echo "✅ Container '$CONTAINER_NAME' exists and is running."

# Check if the network exists
echo "🔍 Checking if the network '$NETWORK_NAME' exists..."
if ! docker network ls --format '{{.Name}}' | grep -w "$NETWORK_NAME"; then
    echo "❌ Network '$NETWORK_NAME' does not exist."
    exit 1
fi
echo "✅ Network '$NETWORK_NAME' exists."

# Connect the container to the network
echo "🔗 Connecting container '$CONTAINER_NAME' to network '$NETWORK_NAME'..."
if docker network connect "$NETWORK_NAME" "$CONTAINER_NAME"; then
    echo "✅ Successfully connected container '$CONTAINER_NAME' to network '$NETWORK_NAME'."
else
    echo "❌ Failed to connect container '$CONTAINER_NAME' to network '$NETWORK_NAME'."
    exit 1
fi