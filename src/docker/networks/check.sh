# Function to check if a Docker network exists
check_network_exists() {
    local NETWORK_NAME=$1

    if [ -z "$NETWORK_NAME" ]; then
        echo "Usage: $0 <network_name>"
        return 0
    fi

    if docker network ls --format '{{.Name}}' | grep -wq "$NETWORK_NAME"; then
        echo "Network '$NETWORK_NAME' already exists."
        return 1
    fi

    return 0
}