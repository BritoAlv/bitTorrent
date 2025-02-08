package serverDocker

import (
	"bittorrent/common"
	"bittorrent/docker/docker_wrapper"
	"github.com/docker/docker/client"
)

var serverLogger = common.NewLogger("DockerServers")
var DockerfilePath ="./docker/server/Dockerfile"
var imageTag = "server:latest"

func SetupServer(containerName string, networkName string, cli *client.Client) error {
	err := docker_wrapper.SetupContainer(DockerfilePath, imageTag, containerName, networkName, cli)
	if err != nil{
		serverLogger.WriteToFileError(err.Error())
		return err
	}
	serverLogger.WriteToFileOK("Success with setting up container for servers")
	return nil
}