check_exist_container() {
    # Check if a container name is provided
    if [ -z "$1" ]; then
        echo "Usage: $0 <container_name>"
        return 0
    fi

    local CONTAINER_NAME=$1

    # Check if a container exists
    if [ "$(docker ps -a -q -f name=^/${CONTAINER_NAME}$)" ]; then
        echo "Container '$CONTAINER_NAME' exists."
        return 1
    else
        echo "Container '$CONTAINER_NAME' does not exist."
        return 0
    fi
}

check_exist_image() {
    if [[ $# -ne 1 ]]; then
        echo "Usage: $0 <image_name>"
        return 0
    fi

    local IMAGE_NAME=$1

    if docker images -q "$IMAGE_NAME" > /dev/null 2>&1; then
        echo "Image '$IMAGE_NAME' exists."
        return 1
    else
        echo "Image '$IMAGE_NAME' does not exist."
        return 0
    fi
}

check_running_container() {
    # Check if a container name is provided
    if [ -z "$1" ]; then
        echo "Usage: $0 <container_name>"
        return 0
    fi

    local CONTAINER_NAME=$1

    # Check if the container exists
    if ! check_exist_container "$CONTAINER_NAME"; then
        return 0
    fi

    # Check if the container is running
    if docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        return 0
    else
        echo "Container '${CONTAINER_NAME}' is running."
        return 1
    fi
}