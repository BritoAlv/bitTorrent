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



# Connect the container to the network
echo "🔗 Connecting container '$CONTAINER_NAME' to network '$NETWORK_NAME'..."
if docker network connect "$NETWORK_NAME" "$CONTAINER_NAME"; then
    echo "✅ Successfully connected container '$CONTAINER_NAME' to network '$NETWORK_NAME'."
else
    echo "❌ Failed to connect container '$CONTAINER_NAME' to network '$NETWORK_NAME'."
    exit 1
fi