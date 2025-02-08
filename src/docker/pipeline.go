package main

import (
	"bittorrent/common"
	clientDocker "bittorrent/docker/client"
	"bittorrent/docker/docker_wrapper"
	routerDocker "bittorrent/docker/router"
	serverDocker "bittorrent/docker/server"
	"fmt"
	"github.com/docker/docker/client"
	"os"
	"strconv"
)

const (
	ClientsNetworkName  = "Clients"
	ServerNetworkName   = "Servers"
	RouterContainerName = "Router"
)

var logger = common.NewLogger("DockerPipeline")

func main() {
	logger.WriteToFileOK("Starting Docker Pipeline")
	/*if len(os.Args) != 3 {
		logger.WriteToFileError("Expecting 2 arguments: <number of clients> and <number of servers>")
		return
	}
	numClients, err := strconv.Atoi(os.Args[1])
	if err != nil {
		logger.WriteToFileError(err.Error())
		return
	}
	numServers, err := strconv.Atoi(os.Args[2])
	if err != nil {
		logger.WriteToFileError(err.Error())
		return
	}*/

	clientWithOpts, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.WriteToFileError(fmt.Sprintf("Failed to create Docker client: %v\n", err))
		os.Exit(1)
	}

	pipeline1(1, 1, clientWithOpts)
}

func pipeline1(numClients int, numServers int, dockerApiClient *client.Client) {
	logger.WriteToFileOK("🚀 Starting the setup...")
	logger.WriteToFileOK("🔧 Creating dockerApiClient network...")
	if err := docker_wrapper.AddNetwork(ClientsNetworkName, dockerApiClient); err != nil {
		logger.WriteToFileError(fmt.Sprintf("Failed to create dockerApiClient network: %v", err))
		return
	}

	logger.WriteToFileOK("🔧 Creating server network...")
	if err := docker_wrapper.AddNetwork(ServerNetworkName, dockerApiClient); err != nil {
		logger.WriteToFileError(fmt.Sprintf("Failed to create server network: %v", err))
		return
	}
	clientsId := make([]string, numClients)
	serversId := make([]string, numServers)

	for i := 1; i <= numClients; i++ {
		clientsId[i-1] = "Client" + strconv.Itoa(i)
		logger.WriteToFileOK(fmt.Sprintf("👤 Adding Client%d to the clients network...", i))
		if err := clientDocker.SetupClient(clientsId[i-1], ClientsNetworkName, dockerApiClient); err != nil {
			logger.WriteToFileError(fmt.Sprintf("Failed to setup Client%d: %v", i, err))
			return
		}
	}

	for i := 1; i <= numServers; i++ {
		serversId[i-1] = "Server" + strconv.Itoa(i)
		logger.WriteToFileOK(fmt.Sprintf("💻 Adding Server%d to the server network...", i))
		if err := serverDocker.SetupServer(serversId[i-1], ServerNetworkName, dockerApiClient); err != nil {
			logger.WriteToFileError(fmt.Sprintf("Failed to setup Server%d: %v", i, err))
			return
		}
	}

	logger.WriteToFileOK("🚦 Starting the router...")
	if err := routerDocker.SetupRouter(RouterContainerName, dockerApiClient); err != nil {
		logger.WriteToFileError(fmt.Sprintf("Failed to setup router: %v", err))
		return
	}

	logger.WriteToFileOK("🔗 Connecting router to dockerApiClient network...")
	if err := docker_wrapper.ConnectContainerNetwork(RouterContainerName, ClientsNetworkName, dockerApiClient); err != nil {
		logger.WriteToFileError(fmt.Sprintf("Failed to connect router to dockerApiClient network: %v", err))
		return
	}

	logger.WriteToFileOK("🔗 Connecting router to server network...")
	if err := docker_wrapper.ConnectContainerNetwork(RouterContainerName, ServerNetworkName, dockerApiClient); err != nil {
		logger.WriteToFileError(fmt.Sprintf("Failed to connect router to server network: %v", err))
		return
	}

	for i := 1; i <= numClients; i++ {
		for j := 1; j <= numServers; j++ {
			logger.WriteToFileOK(fmt.Sprintf("🔍 Testing connection from Client%d to Server%d...", i, j))
			if connected, err := docker_wrapper.TestConnection(clientsId[i-1], serversId[j-1], dockerApiClient); err != nil {
				logger.WriteToFileError(fmt.Sprintf("Failed to test connection from Client%d to Server%d: %v", i, j, err))
				return
			} else if !connected {
				logger.WriteToFileError(fmt.Sprintf("No connection from Client%d to Server%d", i, j))
				return
			}
			logger.WriteToFileOK(fmt.Sprintf("🔍 Testing connection from Server%d to Client%d...", j, i))
			if connected, err := docker_wrapper.TestConnection(clientsId[j-1], serversId[i-1], dockerApiClient); err != nil {
				logger.WriteToFileError(fmt.Sprintf("Failed to test connection from Server%d to Client%d: %v", j, i, err))
				return
			} else if !connected {
				logger.WriteToFileError(fmt.Sprintf("No connection from Server%d to Client%d", j, i))
				return
			}
		}
	}

	for i := 1; i <= numClients; i++ {
		for j := 1; j <= numClients; j++ {
			logger.WriteToFileOK(fmt.Sprintf("🔍 Testing connection from Client%d to Client%d...", i, j))
			if connected, err := docker_wrapper.TestConnection(clientsId[i-1], clientsId[j-1], dockerApiClient); err != nil {
				logger.WriteToFileError(fmt.Sprintf("Failed to test connection from Client%d to Client%d: %v", i, j, err))
				return
			} else if !connected {
				logger.WriteToFileError(fmt.Sprintf("No connection from Client%d to Client%d", i, j))
				return
			}
		}
	}

	for i := 1; i <= numServers; i++ {
		for j := 1; j <= numServers; j++ {
			logger.WriteToFileOK(fmt.Sprintf("🔍 Testing connection from Server%d to Server%d...", i, j))
			if connected, err := docker_wrapper.TestConnection(serversId[i-1], serversId[j-1], dockerApiClient); err != nil {
				logger.WriteToFileError(fmt.Sprintf("Failed to test connection from Server%d to Server%d: %v", i, j, err))
				return
			} else if !connected {
				logger.WriteToFileError(fmt.Sprintf("No connection from Server%d to Server%d", i, j))
				return
			}
		}
	}
	logger.WriteToFileOK("✅ Setup complete!")
}
