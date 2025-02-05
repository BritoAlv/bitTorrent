#!/bin/bash

NETWORK_NAME=$1

if [ -z "$NETWORK_NAME" ]; then
    echo "Usage: $0 <network_name>"
    exit 1
fi

echo "🔄 Checking if the network $NETWORK_NAME exists..."

# Check if the network exists
if ! docker network ls | grep -q "$NETWORK_NAME"; then
    echo "🔄 Network $NETWORK_NAME does not exist. Creating it..."
    if docker network create --driver bridge "$NETWORK_NAME"; then
        echo "✅ Network $NETWORK_NAME created successfully."
    else
        echo "❌ Failed to create network $NETWORK_NAME."
        exit 1
    fi
else
    echo "✅ Network $NETWORK_NAME already exists."
fi

exit 0