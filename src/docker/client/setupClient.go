package clientDocker

import (
	"bittorrent/common"
	"bittorrent/docker/docker_wrapper"
	"github.com/docker/docker/client"
)

var clientLogger = common.NewLogger("DockerClients")
var DockerfilePath = "./docker/client/Dockerfile"
var imageTag = "client:latest"

func SetupClient(containerName string, networkName string, cli *client.Client) error {
	err := docker_wrapper.SetupContainer(DockerfilePath, imageTag, containerName, networkName, cli)
	if err != nil {
		clientLogger.WriteToFileError(err.Error())
		return err
	}
	clientLogger.WriteToFileOK("Success with setting up container for clients")
	return nil
}
