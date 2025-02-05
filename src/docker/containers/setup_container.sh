NETWORK_SCRIPTS_PATH=./docker/networks/

if [ "$#" -lt 3 ] || [ "$#" -gt 4 ]; then
    echo "Usage: $0 DOCKER_FILE_PATH IMAGE_TAG CONTAINER_NAME [NETWORK]"
    exit 1
fi

DOCKER_FILE_PATH=$1
IMAGE_TAG=$2
CONTAINER_NAME=$3

echo "🔨 Building Docker image..."
if sh build_image.sh $IMAGE_TAG $DOCKER_FILE_PATH then 
    exit 1
fi
echo "✅ Docker image built."

echo "🚀 Starting Docker container..."
if sh start_container.sh $IMAGE_TAG $CONTAINER_NAME then
    exit 1
fi
echo "✅ Docker container started."

echo "🔍 Accessing Docker container..."
if sh access_container.sh $CONTAINER_NAME then
    exit 1
fi
echo "✅ Accessed Docker container."

if [ "$#" -eq 4 ]; then
    NETWORK=$4
    echo "🔗 Connecting container to network..."
    if sh ${NETWORK_SCRIPTS_PATH}connect_container_network.sh $CONTAINER_NAME $NETWORK then
        exit 1
    fi
    echo "✅ Container connected to network."
fi

exit 0