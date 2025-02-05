# Check if exactly one argument is provided.
if [ "$#" -ne 1 ]; then
    echo "❌ Usage: $0 <container_name>"
    exit 1
fi
echo "✅ Container name provided: $CONTAINER_NAME"

# Check if the container is running.
if sh check_running_container.sh $1; then
    exit 1
fi

if docker exec --rm -it $CONTAINER_NAME /bin/sh; then
    exit 1
fi

exit 0