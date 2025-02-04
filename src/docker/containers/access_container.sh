# Name of the container
CONTAINER_NAME=$1

# Check if the container name is provided
if [ -z "$CONTAINER_NAME" ]; then
    echo "❌ Usage: $0 <container_name>"
    exit 1
else
    echo "✅ Container name provided: $CONTAINER_NAME"
fi

# Check if the container is running
if [ "$(docker ps -q -f name=$CONTAINER_NAME)" ]; then
    echo "✅ Container $CONTAINER_NAME is running."
    echo "Accessing the container $CONTAINER_NAME..."
    docker exec -it $CONTAINER_NAME /bin/sh
else
    echo "❌ Container $CONTAINER_NAME is not running."
    exit 1
fi