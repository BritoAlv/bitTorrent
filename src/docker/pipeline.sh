# This pipeline script is used to create a network of clients and servers, all the clients are in the same network and all the servers are in the same network (but other).

# Exit immediately if a command exits with a non-zero status.
set -e

NETWORK_SCRIPTS_PATH=./docker/networks/
CLIENTS_PATH=./docker/client/
SERVERS_PATH=./docker/server/
ROUTER_PATH=./docker/router/
CLIENTS_NETWORK_NAME=Clients
SERVER_NETWORK_NAME=Servers
ROUTER_CONTAINER_NAME=Router


# Check if the correct number of arguments are provided.
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <number_of_clients> <number_of_servers>"
    exit 1
fi

NUM_CLIENTS=$1
NUM_SERVERS=$2

echo "🚀 Starting the setup..."

echo "🔧 Creating client network..."
sh ${NETWORK_SCRIPTS_PATH}add_network.sh $CLIENTS_NETWORK_NAME

echo "🔧 Creating server network..."
sh ${NETWORK_SCRIPTS_PATH}add_network.sh $SERVER_NETWORK_NAME

for i in $(seq 1 "$NUM_CLIENTS"); do
    echo "👤 Adding Client${i} to the clients network..."
    sh ${CLIENTS_PATH}run.sh Client"${i}" $CLIENTS_NETWORK_NAME
done

# Add servers to the servers network
for i in $(seq 1 "$NUM_SERVERS"); do
    echo "💻 Adding Server${i} to the server network..."
    sh ${SERVERS_PATH}run.sh Server"${i}" $SERVER_NETWORK_NAME
done

echo "🚦 Starting the router..."
sh ${ROUTER_PATH}run.sh $ROUTER_CONTAINER_NAME


echo "🔗 Connecting router to client network..."
sh ${NETWORK_SCRIPTS_PATH}connect_container_network.sh $ROUTER_CONTAINER_NAME $CLIENTS_NETWORK_NAME

echo "🔗 Connecting router to server network..."
sh ${NETWORK_SCRIPTS_PATH}connect_container_network.sh $ROUTER_CONTAINER_NAME $SERVER_NETWORK_NAME

# Test connection between the clients and the servers.
for i in $(seq 1 "$NUM_CLIENTS"); do
    for j in $(seq 1 "$NUM_SERVERS"); do
        echo "🔍 Testing connection from Client${i} to Server${j}..."
        sh ${NETWORK_SCRIPTS_PATH}/test_connection.sh Client"${i}" Server"${j}"
        echo "🔍 Testing connection from Server${j} to Client${i}..."
        sh ${NETWORK_SCRIPTS_PATH}/test_connection.sh Server"${j}" Client"${i}"
    done
done

for i in $(seq 1 "$NUM_CLIENTS"); do
    for j in $(seq 1 "$NUM_CLIENTS"); do
        echo "🔍 Testing connection from Client${i} to Client${j}..."
        sh ${NETWORK_SCRIPTS_PATH}test_connection.sh Client"${i}" Client"${j}"
    done
done

for i in $(seq 1 "$NUM_SERVERS"); do
    for j in $(seq 1 "$NUM_SERVERS"); do
        echo "🔍 Testing connection from Server${i} to Server${j}..."
        sh test_connection.sh Server"${i}" Server"${j}"
    done
done

echo "✅ Setup complete!"