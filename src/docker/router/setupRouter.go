package routerDocker

import (
	"bittorrent/common"
	"bittorrent/docker/docker_wrapper"
	"github.com/docker/docker/client"
)

var routerLogger = common.NewLogger("DockerRouter")
var DockerfilePath = "./docker/router/Dockerfile"
var imageTag = "router:latest"

func SetupRouter(containerName string, cli *client.Client) error {
	err := docker_wrapper.SetupContainer(DockerfilePath, imageTag, containerName, "", cli)
	if err != nil {
		routerLogger.WriteToFileError(err.Error())
		return err
	}
	routerLogger.WriteToFileOK("Success with setting up container for routers")
	return nil
}
