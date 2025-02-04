#!/bin/bash

# Check if network name is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <network_name>"
    exit 1
fi

NETWORK_NAME=$1
# Remove the Docker network
if docker network rm "$NETWORK_NAME"; then
    echo "✅ Network '$NETWORK_NAME' removed successfully."
else
    echo "❌ Failed to remove network '$NETWORK_NAME'."
fi