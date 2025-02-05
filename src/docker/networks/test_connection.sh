#!/bin/bash

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
if docker exec "$SOURCE_CONTAINER" ping -c 4 "$TARGET_CONTAINER"; then
    echo "❌ Connection failed between $SOURCE_CONTAINER and $TARGET_CONTAINER"
    exit 1
else
    echo "✅ Connection successful between $SOURCE_CONTAINER and $TARGET_CONTAINER"
    exit 0
fi