# Check if the correct number of arguments is provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <image_tag> <container_name>"
    exit 1
fi

# Assign arguments to variables
IMAGE_TAG=$1
CONTAINER_NAME=$2

# Check if the Docker image exists
if sh check_image.sh "$IMAGE_TAG"; then
    echo "❌ Docker image with tag '$IMAGE_TAG' does not exist."
    exit 1
fi
echo "✅ Docker image with tag '$IMAGE_TAG' found."

# Check if a container with the same name is already running
if sh check_running_container.sh "$CONTAINER_NAME"; then
    echo "❌ A container with the name '$CONTAINER_NAME' is already running."
    exit 1
fi
echo "✅ No container with the name '$CONTAINER_NAME' is running."

# Start a docker container and open a shell
# rm: Automatically remove the container once it stops (useful for cleanup).
# it: Run the container in interactive mode.
# privileged: Give the container full access to the host system.
# name: Assign a name to the container.
if docker run --rm --privileged --name "$CONTAINER_NAME" "$TAG" then 
    exit 1
fi

exit 0