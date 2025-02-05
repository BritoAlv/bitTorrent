# Commands for building s docker image given a tag and a dockerfile.
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <tag> <docker_file_path>"
    exit 1
fi

TAG=$1
DOCKER_FILE_PATH=$2

echo "✅ Tag provided: $TAG"
echo "✅ Docker file path provided: $DOCKER_FILE_PATH"

# Build the docker image
# t: Assign a tag to the image, but its like name:tag.

if sh check_image.sh $TAG; then
    echo "Image '$TAG' exists."
    exit 1
fi

if docker build -f $DOCKER_FILE_PATH -t $TAG . ; then
    exit 1
fi

exit 0