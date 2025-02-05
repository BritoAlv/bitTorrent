#!/bin/bash
set -e
# Check if two arguments are provided
if [ "$#" -ne 2 ]; then
    echo "❌ Usage: $0 <source_container> <target_container>"
    exit 1
fi

SOURCE_CONTAINER=$1
TARGET_CONTAINER=$2

# Execute ping command from source container to target container
echo "🔄 Pinging $TARGET_CONTAINER from $SOURCE_CONTAINER..."

# Check the exit status of the ping command
docker exec "$SOURCE_CONTAINER" ping -c 4 "$TARGET_CONTAINER"
echo "✅ Connection successful between $SOURCE_CONTAINER and $TARGET_CONTAINER"