set -e

NETWORK_SCRIPTS_PATH=./docker/networks/

if [ "$#" -lt 3 ] || [ "$#" -gt 4 ]; then
    echo "Usage: $0 DOCKER_FILE_PATH IMAGE_TAG CONTAINER_NAME [NETWORK]"
    exit 1
fi

DOCKER_FILE_PATH=$1
IMAGE_TAG=$2
CONTAINER_NAME=$3

echo "🔨 Building Docker image..."
sh build_image.sh $IMAGE_TAG $DOCKER_FILE_PATH
echo "✅ Docker image built."

"🚀 Starting Docker container..."
sh start_container.sh $IMAGE_TAG $CONTAINER_NAME
echo "✅ Docker container started."

echo "🔍 Accessing Docker container..."
sh access_container.sh $CONTAINER_NAME
echo "✅ Accessed Docker container."

if [ "$#" -eq 4 ]; then
    NETWORK=$4
    echo "🔗 Connecting container to network..."
    sh ${NETWORK_SCRIPTS_PATH}connect_container_network.sh $CONTAINER_NAME $NETWORK    
    echo "✅ Container connected to network."
fi