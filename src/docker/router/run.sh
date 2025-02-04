CONTAINER_SCRIPTS_PATH=../containers/
IMAGE_TAG=router:latest
CWD=$(pwd)

if [ $# -ne 1 ]; then
    echo "Usage: $0 <container_name> "
    exit 1
fi

CONTAINER_NAME=$1
sh ${CONTAINER_SCRIPTS_PATH}setup_container.sh $CWD $IMAGE_TAG $CONTAINER_NAME