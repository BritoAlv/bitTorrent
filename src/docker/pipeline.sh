# This pipeline script is used to create a network of clients and servers, all the clients are in the same network and all the servers are in the same network (but other).

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

# Create client and server networks using your scripts.
echo "🔧 Creating client network..."
if sh ${NETWORK_SCRIPTS_PATH}add_network.sh $CLIENTS_NETWORK_NAME ; then
   exit 1
fi

echo "🔧 Creating server network..."
if sh ${NETWORK_SCRIPTS_PATH}add_network.sh $SERVER_NETWORK_NAME ; then
   exit 1
fi

# Add clients to the clients network
for i in $(seq 1 "$NUM_CLIENTS"); do
    echo "👤 Adding Client${i} to the clients network..."
    if sh ${CLIENTS_PATH}run.sh Client"${i}" $CLIENTS_NETWORK_NAME ; then
        exit 1
    fi
done

# Add servers to the servers network
for i in $(seq 1 "$NUM_SERVERS"); do
    echo "💻 Adding Server${i} to the server network..."
    if sh ${SERVERS_PATH}run.sh Server"${i}" $SERVER_NETWORK_NAME ; then
        exit 1
    fi
done

echo "🚦 Starting the router..."

if sh ${ROUTER_PATH}run.sh $ROUTER_CONTAINER_NAME ; then
    exit 1
fi

echo "🔗 Connecting router to client network..."
if sh ${NETWORK_SCRIPTS_PATH}connect_container_network.sh $ROUTER_CONTAINER_NAME $CLIENTS_NETWORK_NAME ; then
    exit 1
fi

echo "🔗 Connecting router to server network..."
if sh ${NETWORK_SCRIPTS_PATH}connect_container_network.sh $ROUTER_CONTAINER_NAME $SERVER_NETWORK_NAME ; then
    exit 1
fi

# Test connection between the clients and the servers.
for i in $(seq 1 "$NUM_CLIENTS"); do
    for j in $(seq 1 "$NUM_SERVERS"); do
        echo "🔍 Testing connection from Client${i} to Server${j}..."
        if sh ${NETWORK_SCRIPTS_PATH}/test_connection.sh Client"${i}" Server"${j}" ; then
            exit 1
        fi
        echo "🔍 Testing connection from Server${j} to Client${i}..."
        if sh ${NETWORK_SCRIPTS_PATH}/test_connection.sh Server"${j}" Client"${i}" ; then
            exit 1
        fi
    done
done

for i in $(seq 1 "$NUM_CLIENTS"); do
    for j in $(seq 1 "$NUM_CLIENTS"); do
        echo "🔍 Testing connection from Client${i} to Client${j}..."
        if sh ${NETWORK_SCRIPTS_PATH}test_connection.sh Client"${i}" Client"${j}" ; then
            exit 1
        fi
    done
done

for i in $(seq 1 "$NUM_SERVERS"); do
    for j in $(seq 1 "$NUM_SERVERS"); do
        echo "🔍 Testing connection from Server${i} to Server${j}..."
        if sh test_connection.sh Server"${i}" Server"${j}" ; then
            exit 1
        fi
    done
done

echo "✅ Setup complete!"
exit 0