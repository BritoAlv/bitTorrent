#!/bin/bash
set -e
# Check if network name is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <network_name>"
    exit 1
fi

NETWORK_NAME=$1

# Remove the Docker network
docker network rm "$NETWORK_NAME"
echo "✅ Network '$NETWORK_NAME' removed successfully."