CONTAINER_SCRIPTS_PATH=./docker/containers/
DOCKER_FILE_PATH=./docker/client/Dockerfile
IMAGE_TAG=router:latest

if [ $# -ne 1 ]; then
    echo "Usage: $0 <container_name> "
    exit 1
fi

CONTAINER_NAME=$1

if sh ${CONTAINER_SCRIPTS_PATH}setup_container.sh $DOCKER_FILE_PATH $IMAGE_TAG $CONTAINER_NAME; then 
    exit 1
fi

exit 0