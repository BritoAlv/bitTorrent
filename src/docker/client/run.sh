CONTAINER_SCRIPTS_PATH=../containers/
IMAGE_TAG=client:latest
CWD=$(pwd)

if [ $# -ne 2 ]; then
    echo "Usage: $0 <container_name> <network_name>"
    exit 1
fi

CONTAINER_NAME=$1
NETWORK_NAME=$2
sh ${CONTAINER_SCRIPTS_PATH}setup_container.sh $CWD $IMAGE_TAG $CONTAINER_NAME $NETWORK_NAME 