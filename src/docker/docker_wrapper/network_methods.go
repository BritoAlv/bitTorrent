package docker_wrapper

import (
	"bittorrent/common"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"os/exec"
)

var networkLogger = common.NewLogger("DockerNetwork")

func DisconnectContainerNetwork(containerName string, networkName string, cli *client.Client) error {
	networkLogger.WriteToFileOK(fmt.Sprintf("DisconnectContainerNetwork(containerName=%s, networkName=%s)", containerName, networkName))
	// Check if the container is running
	isRunning, err := CheckContainerIsRunning(containerName, cli)
	if err != nil {
		networkLogger.WriteToFileError(fmt.Sprintf("Failed to check if container %s is running: %v", containerName, err))
		return err
	}
	if !isRunning {
		networkLogger.WriteToFileError(fmt.Sprintf("Container %s is not running", containerName))
		return fmt.Errorf("container %s is not running", containerName)
	}

	// Check if the network exists
	networkExists, err := ExistNetwork(networkName, cli)
	if err != nil {
		networkLogger.WriteToFileError(fmt.Sprintf("Failed to check if network %s exists: %v", networkName, err))
		return err
	}
	if !networkExists {
		networkLogger.WriteToFileError(fmt.Sprintf("Network %s does not exist", networkName))
		return fmt.Errorf("network %s does not exist", networkName)
	}

	// Disconnect the container from the network
	err = cli.NetworkDisconnect(context.Background(), networkName, containerName, false)
	if err != nil {
		networkLogger.WriteToFileError(fmt.Sprintf("Failed to disconnect container %s from network %s: %v", containerName, networkName, err))
		return err
	}

	networkLogger.WriteToFileOK(fmt.Sprintf("Container %s disconnected from network %s successfully", containerName, networkName))
	return nil
}

func ConnectContainerNetwork(containerName string, networkName string, cli *client.Client) error {
	networkLogger.WriteToFileOK(fmt.Sprintf("ConnectContainerNetwork(containerName=%s, networkName=%s)", containerName, networkName))
	// Check if the container is running
	exist, err := CheckExistContainer(containerName, cli)
	if err != nil {
		networkLogger.WriteToFileError(fmt.Sprintf("Failed to check if container %s is running: %v", containerName, err))
		return err
	}
	if !exist {
		networkLogger.WriteToFileError(fmt.Sprintf("Container %s is not running", containerName))
		return fmt.Errorf("container %s is not running", containerName)
	}

	// Check if the network exists
	networkExists, err := ExistNetwork(networkName, cli)
	if err != nil {
		networkLogger.WriteToFileError(fmt.Sprintf("Failed to check if network %s exists: %v", networkName, err))
		return err
	}
	if !networkExists {
		networkLogger.WriteToFileError(fmt.Sprintf("Network %s does not exist", networkName))
		return fmt.Errorf("network %s does not exist", networkName)
	}

	// Connect the container to the network
	err = cli.NetworkConnect(context.Background(), networkName, containerName, nil)
	if err != nil {
		networkLogger.WriteToFileError(fmt.Sprintf("Failed to connect container %s to network %s: %v", containerName, networkName, err))
		return err
	}

	networkLogger.WriteToFileOK(fmt.Sprintf("Container %s connected to network %s successfully", containerName, networkName))
	return nil
}

func ExistNetwork(networkName string, cli *client.Client) (bool, error) {
	networkLogger.WriteToFileOK(fmt.Sprintf("ExistNetwork(networkName=%s)", networkName))
	networks, err := cli.NetworkList(context.Background(), network.ListOptions{})
	if err != nil {
		networkLogger.WriteToFileError(fmt.Sprintf("Something wrong happened when checking if the network %s exists: %s", networkName, err.Error()))
		return false, err
	}
	for _, networkFound := range networks {
		if networkFound.Name == networkName {
			return true, nil
		}
	}
	return false, nil
}

func AddNetwork(networkName string, cli *client.Client) error {
	networkLogger.WriteToFileOK(fmt.Sprintf("AddNetwork(networkName=%s)", networkName))
	result, err := ExistNetwork(networkName, cli)
	if err != nil {
		networkLogger.WriteToFileError(fmt.Sprintf("Failed to add network %s: %s", networkName, err.Error()))
		return err
	}
	if result {
		networkLogger.WriteToFileOK(fmt.Sprintf("Network %s already exists", networkName))
		return nil
	}
	_, err = cli.NetworkCreate(context.Background(), networkName, network.CreateOptions{
		Driver: "bridge",
		IPAM: &network.IPAM{
			Driver: "default",
		},
	})
	if err != nil {
		networkLogger.WriteToFileError(fmt.Sprintf("Failed to create network %s: %v", networkName, err))
		return err
	}
	networkLogger.WriteToFileOK(fmt.Sprintf("Network %s created successfully", networkName))
	return nil
}

func RemoveNetwork(networkName string, cli *client.Client) error {
	networkLogger.WriteToFileOK(fmt.Sprintf("RemoveNetwork(networkName=%s)", networkName))
	result, err := ExistNetwork(networkName, cli)
	if err != nil {
		networkLogger.WriteToFileError(fmt.Sprintf("Failed to check if network %s exists: %s", networkName, err.Error()))
		return err
	}
	if !result {
		networkLogger.WriteToFileOK(fmt.Sprintf("Network %s does not exist", networkName))
		return nil
	}
	err = cli.NetworkRemove(context.Background(), networkName)
	if err != nil {
		networkLogger.WriteToFileError(fmt.Sprintf("Failed to remove network %s: %v", networkName, err))
		return err
	}
	networkLogger.WriteToFileOK(fmt.Sprintf("Network %s removed successfully", networkName))
	return nil
}

func TestConnection(containerOne string, containerTwo string, cli *client.Client) (bool, error) {
	networkLogger.WriteToFileOK(fmt.Sprintf("TestConnection(containerOne=%s, containerTwo=%s)", containerOne, containerTwo))
	// Check if containerOne is running
	isRunningOne, err := CheckContainerIsRunning(containerOne, cli)
	if err != nil {
		networkLogger.WriteToFileError(fmt.Sprintf("Failed to check if container %s is running: %v", containerOne, err))
		return false, err
	}
	if !isRunningOne {
		networkLogger.WriteToFileError(fmt.Sprintf("Container %s is not running", containerOne))
		return false, nil
	}

	// Check if containerTwo is running
	isRunningTwo, err := CheckContainerIsRunning(containerTwo, cli)
	if err != nil {
		networkLogger.WriteToFileError(fmt.Sprintf("Failed to check if container %s is running: %v", containerTwo, err))
		return false, err
	}
	if !isRunningTwo {
		networkLogger.WriteToFileError(fmt.Sprintf("Container %s is not running", containerTwo))
		return false, nil
	}
	// Execute ping command from containerOne to containerTwo
	cmd := exec.Command("docker", "exec", containerOne, "ping", "-c", "4", containerTwo)
	output, err := cmd.CombinedOutput()
	if err != nil {
		networkLogger.WriteToFileError(fmt.Sprintf("Failed to ping %s from %s: %v\nOutput: %s", containerTwo, containerOne, err, string(output)))
		return false, err
	}
	networkLogger.WriteToFileOK(fmt.Sprintf("Connection successful between %s and %s\nOutput: %s", containerOne, containerTwo, string(output)))
	return true, nil
}
