#!/bin/bash
set -e
NETWORK_NAME=$1

source ./docker/networks/check.sh

if [ -z "$NETWORK_NAME" ]; then
    echo "Usage: $0 <network_name>"
    exit 1
fi

echo "🔄 Checking if the network $NETWORK_NAME exists..."
if check_network_exists $NETWORK_NAME ; then
    echo "⚠️ Network $NETWORK_NAME already exists."
    exit 0
fi

docker network create --driver bridge "$NETWORK_NAME"
echo "✅ Network $NETWORK_NAME created successfully."