set -e
CONTAINER_SCRIPTS_PATH=./docker/containers/
NETWORK_SCRIPTS_PATH=./docker/networks/
IMAGE_TAG=client:latest
DOCKERFILE_PATH=./docker/client/Dockerfile

if [ $# -ne 2 ]; then
    echo "Usage: $0 <container_name> <network_name>"
    exit 1
fi

CONTAINER_NAME=$1
NETWORK_NAME=$2

sh ${CONTAINER_SCRIPTS_PATH}setup_container.sh $DOCKER_FILE_PATH $IMAGE_TAG $CONTAINER_NAME $NETWORK_NAME