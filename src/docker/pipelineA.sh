#!/bin/bash
NETWORK_SCRIPTS_PATH=./networks/
CLIENTS_NETWORK_NAME=Clients
SERVER_NETWORK_NAME=Servers
ROUTER_CONTAINER_NAME=Router
# Check if the correct number of arguments are provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <number_of_clients> <number_of_servers>"
    exit 1
fi

NUM_CLIENTS=$1
NUM_SERVERS=$2

# Create client and server networks using your scripts
sh ${NETWORK_SCRIPTS_PATH}add_network.sh $CLIENTS_NETWORK_NAME
sh ${NETWORK_SCRIPTS_PATH}add_network.sh $SERVER_NETWORK_NAME

# Add clients to the clients network
for i in $(seq 1 $NUM_CLIENTS); do
    sh ./client/run.sh Client${i} $CLIENTS_NETWORK_NAME
done

# Add servers to the servers network
for i in $(seq 1 $NUM_SERVERS); do
    sh ./server/run.sh Server${i} $SERVER_NETWORK_NAME
done

sh ./router/run.sh $ROUTER_CONTAINER_NAME

sh ${NETWORK_SCRIPTS_PATH}connect_container_network.sh $ROUTER_CONTAINER_NAME $CLIENTS_NETWORK_NAME
sh ${NETWORK_SCRIPTS_PATH}connect_container_network.sh $ROUTER_CONTAINER_NAME $SERVERS_NETWORK_NAME