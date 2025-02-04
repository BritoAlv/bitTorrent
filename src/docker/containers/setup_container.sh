if [ "$#" -lt 3 ] || [ "$#" -gt 4 ]; then
    echo "Usage: $0 DOCKER_FILE_PATH IMAGE_TAG CONTAINER_NAME [NETWORK]"
    exit 1
fi

DOCKER_FILE_PATH=$1
IMAGE_TAG=$2
CONTAINER_NAME=$3
NETWORK_SCRIPTS_PATH=../networks/ 

sh build_image.sh $IMAGE_TAG $CWD
sh start_container.sh $IMAGE_TAG $CONTAINER_NAME
sh access_container.sh $CONTAINER_NAME

if [ "$#" -eq 4 ]; then
    NETWORK=$4
    sh ${NETWORK_SCRIPTS_PATH}connect_container_network.sh $CONTAINER_NAME $NETWORK
fi