CONTAINERS_SCRIPT_PATH=./docker/containers/

# Check if the correct number of arguments is provided
if [ "$#" -ne 2 ]; then
    echo "❌ Usage: $0 <container_name_or_id> <network_name>"
    exit 1
fi

CONTAINER_NAME_OR_ID=$1
NETWORK_NAME=$2

if sh ${CONTAINERS_SCRIPT_PATH}check_container.sh "$CONTAINER_NAME_OR_ID"; then
    echo "❌ Container $CONTAINER_NAME_OR_ID does not exist."
    exit 1
else
    echo "✅ Container $CONTAINER_NAME_OR_ID exists."
fi 

# Check if the network exists
if sh check_network.sh "$NETWORK_NAME"; then
    echo "❌ Network $NETWORK_NAME does not exist."
    exit 1
else
    echo "✅ Network $NETWORK_NAME exists."
fi


echo "🔍 Checking if the container is connected to the network..."

# Disconnect the container from the network
if docker network disconnect "$NETWORK_NAME" "$CONTAINER_NAME_OR_ID"; then
    echo "❌ Failed to disconnect $CONTAINER_NAME_OR_ID from $NETWORK_NAME"
    exit 1
else
    echo "✅ Successfully disconnected $CONTAINER_NAME_OR_ID from $NETWORK_NAME"
fi

exit 0