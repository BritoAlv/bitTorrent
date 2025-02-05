set -e
source ./docker/containers/check.sh
# Check if exactly one argument is provided.
if [ "$#" -ne 1 ]; then
    echo "❌ Usage: $0 <container_name>"
    exit 1
fi

echo "✅ Container name provided: $CONTAINER_NAME"
# Check if the container is running.
if !check_running_container $1 ; then
    docker exec --rm -it $CONTAINER_NAME /bin/sh
else 
    echo "Container is already running."
fi