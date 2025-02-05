CONTAINER_SCRIPTS_PATH=./docker/containers
IMAGE_TAG=server:latest
DOCKERFILE_PATH=./docker/server/Dockerfile

if [ $# -ne 2 ]; then
    echo "Usage: $0 <container_name> <network_name>"
    exit 1
fi

CONTAINER_NAME=$1
NETWORK_NAME=$2

if sh ${CONTAINER_SCRIPTS_PATH}setup_container.sh $DOCKERFILE_PATH $IMAGE_TAG $CONTAINER_NAME $NETWORK_NAME; then
    exit 1
fi

exit 0