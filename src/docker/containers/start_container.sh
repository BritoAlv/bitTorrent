set -e
source ./docher/containers/check.sh
# Check if the correct number of arguments is provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <image_tag> <container_name>"
    exit 1
fi

# Assign arguments to variables
IMAGE_TAG=$1
CONTAINER_NAME=$2

# Check if the Docker image exists
if ! check_exist_image $IMAGE_TAG; then 
    echo "❌ Docker image with tag '$IMAGE_TAG' not found."
    exit 1
fi
echo "✅ Docker image with tag '$IMAGE_TAG' found."

# Check if a container with the same name is already running
if check_running_container.sh $CONTAINER_NAME ; then 
    echo "❌ A container with the name '$CONTAINER_NAME' is already running."
else 
    echo "✅ No container with the name '$CONTAINER_NAME' is running."
    docker run --rm --privileged --name "$CONTAINER_NAME" "$TAG"
fi