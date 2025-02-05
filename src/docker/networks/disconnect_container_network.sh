set -e
CONTAINERS_SCRIPT_PATH=./docker/containers/

source ./docker/networks/check.sh
source ./docker/containers/check.sh

# Check if the correct number of arguments is provided
if [ "$#" -ne 2 ]; then
    echo "❌ Usage: $0 <container_name> <network_name>"
    exit 1
fi

CONTAINER_NAME=$1
NETWORK_NAME=$2

# Connect the container to the network
if ! check_exist_container $CONTAINER_NAME; then
    echo "❌ Container '$CONTAINER_NAME' not found."
    exit 1
fi

if ! check_network_exists $NETWORK_NAME; then
    echo "❌ Network '$NETWORK_NAME' not found."
    exit 1
fi

# Disconnect the container from the network
docker network disconnect "$NETWORK_NAME" "$CONTAINER_NAME"
echo "✅ Successfully disconnected $CONTAINER_NAME from $NETWORK_NAME"