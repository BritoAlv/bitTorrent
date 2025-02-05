#!/bin/bash
if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <image_name>"
    exit 1
fi

IMAGE_NAME=$1

if docker images -q "$IMAGE_NAME" > /dev/null 2>&1; then
    echo "Image '$IMAGE_NAME' exists."
    exit 0
else
    echo "Image '$IMAGE_NAME' does not exist."
    exit 1
fi