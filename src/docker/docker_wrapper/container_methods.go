package docker_wrapper

import (
	"archive/tar"
	"bittorrent/common"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

var containerLogger = common.NewLogger("DockerContainerLoggers")

func SetupContainer(dockerFilePath, imageTag, containerName, networkName string, cli *client.Client) error {
	containerLogger.WriteToFileOK(fmt.Sprintf("SetupContainer(dockerFilePath=%s, imageTag=%s, containerName=%s, networkName=%s)", dockerFilePath, imageTag, containerName, networkName))
	if err := BuildImage(imageTag, dockerFilePath, cli); err != nil {
		return fmt.Errorf("failed to build Docker image: %v", err)
	}

	if err := CreateContainer(imageTag, containerName, cli); err != nil {
		return fmt.Errorf("failed to create Docker container: %v", err)
	}

	if networkName != "" {
		if err := ConnectContainerNetwork(containerName, networkName, cli); err != nil {
			return fmt.Errorf("failed to connect container to network: %v", err)
		}
	}

	if err := RunDockerContainer(containerName, cli); err != nil {
		return fmt.Errorf("failed to run Docker container: %v", err)
	}

	if err := AccessContainer(containerName, cli); err != nil {
		return fmt.Errorf("failed to access Docker container: %v", err)
	}

	return nil
}

func RunDockerContainer(containerName string, cli *client.Client) error {
	containerLogger.WriteToFileOK(fmt.Sprintf("RunDockerContainer(containerName=%s)", containerName))
	ctx := context.Background()

	containerInspect, err := cli.ContainerInspect(ctx, containerName)
	if err != nil {
		containerLogger.WriteToFileError(fmt.Sprintf("Failed to inspect container %s: %v", containerName, err))
		return err
	}

	err = cli.ContainerStart(ctx, containerInspect.ID, container.StartOptions{})
	if err != nil {
		containerLogger.WriteToFileError(fmt.Sprintf("Failed to start container %s: %v", containerName, err))
		return err
	}
	containerLogger.WriteToFileOK(fmt.Sprintf("Container %s started successfully", containerName))
	return nil
}

func CreateContainer(imageTag string, containerName string, cli *client.Client) error {
	containerLogger.WriteToFileOK(fmt.Sprintf("CreateContainer(imageTag=%s, containerName=%s)", imageTag, containerName))
	exists, err := CheckExistContainer(containerName, cli)
	if err != nil {
		containerLogger.WriteToFileError(fmt.Sprintf("Failed to check if container %s exists: %v", containerName, err))
		return err
	}
	if exists {
		containerLogger.WriteToFileOK(fmt.Sprintf("Container %s already exists", containerName))
		return nil
	}

	_, err = cli.ContainerCreate(context.Background(), &container.Config{
		Image: imageTag,
	}, &container.HostConfig{
		Privileged: true, // Enable privileged mode
	}, nil, nil, containerName)
	if err != nil {
		containerLogger.WriteToFileError(fmt.Sprintf("Failed to create container %s: %v", containerName, err))
		return err
	}

	containerLogger.WriteToFileOK(fmt.Sprintf("Container %s created successfully", containerName))
	return nil
}

func AccessContainer(containerName string, cli *client.Client) error {
	containerLogger.WriteToFileOK(fmt.Sprintf("AccessContainer(containerName=%s)", containerName))
	isRunning, err := CheckContainerIsRunning(containerName, cli)
	if err != nil {
		containerLogger.WriteToFileError(fmt.Sprintf("Failed to check if container %s is running: %v", containerName, err))
		return err
	}
	if !isRunning {
		containerLogger.WriteToFileError(fmt.Sprintf("Container %s is not running", containerName))
		return fmt.Errorf("container %s is not running", containerName)
	}
	cmd := exec.Command("docker", "exec", "-it", containerName, "/bin/sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		containerLogger.WriteToFileError(fmt.Sprintf("Failed to open container %s in terminal: %v", containerName, err))
		return err
	}

	containerLogger.WriteToFileOK(fmt.Sprintf("Container %s opened in terminal successfully", containerName))
	return nil
}

func BuildImage(imageTag string, dockerFilePath string, cli *client.Client) error {
	containerLogger.WriteToFileOK(fmt.Sprintf("BuildImage(imageTag=%s, dockerFilePath=%s)", imageTag, dockerFilePath))
	dockerFileReader, err := os.Open(dockerFilePath)
	if err != nil {
		containerLogger.WriteToFileError(fmt.Sprintf("Failed to open Dockerfile %s: %v", dockerFilePath, err))
		return err
	}
	defer dockerFileReader.Close()

	// Get directory of Dockerfile to use as build context
	dockerFileDir := filepath.Dir("./")

	// Tar the directory (Docker expects a tarred context)
	tarReader, err := TarDirectory(dockerFileDir)
	if err != nil {
		containerLogger.WriteToFileError(fmt.Sprintf("Failed to create tar for directory %s: %v", dockerFileDir, err))
		return err
	}

	buildOptions := types.ImageBuildOptions{
		Tags:       []string{imageTag},
		Dockerfile: dockerFilePath,
		Remove:     true,
	}

	buildResponse, err := cli.ImageBuild(context.Background(), tarReader, buildOptions)
	if err != nil {
		containerLogger.WriteToFileError(fmt.Sprintf("Failed to build image %s: %v", imageTag, err))
		return err
	}
	// Stream logs properly
	scanner := bufio.NewScanner(buildResponse.Body)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		containerLogger.WriteToFileOK(line) // Log each line
	}

	if err := scanner.Err(); err != nil {
		containerLogger.WriteToFileError(fmt.Sprintf("Error reading build response: %v", err))
		return err
	}

	containerLogger.WriteToFileOK(fmt.Sprintf("Image %s built successfully", imageTag))
	return nil
}

func CheckExistContainer(containerName string, cli *client.Client) (bool, error) {
	containerLogger.WriteToFileOK(fmt.Sprintf("CheckExistContainer(containerName=%s)", containerName))
	_, err := cli.ContainerInspect(context.Background(), containerName)
	if err != nil {
		if client.IsErrNotFound(err) {
			containerLogger.WriteToFileError(fmt.Sprintf("Container %s does not exist", containerName))
			return false, nil
		}
		containerLogger.WriteToFileError(fmt.Sprintf("Failed to inspect container %s: %v", containerName, err))
		return false, err
	}
	containerLogger.WriteToFileOK(fmt.Sprintf("Container %s exists", containerName))
	return true, nil
}

func CheckContainerIsRunning(containerName string, cli *client.Client) (bool, error) {
	containerLogger.WriteToFileOK(fmt.Sprintf("CheckContainerIsRunning(containerName=%s)", containerName))
	containerInspect, err := cli.ContainerInspect(context.Background(), containerName)
	if err != nil {
		if client.IsErrNotFound(err) {
			containerLogger.WriteToFileError(fmt.Sprintf("Container %s does not exist", containerName))
			return false, nil
		}
		containerLogger.WriteToFileError(fmt.Sprintf("Failed to inspect containerInspect %s: %v", containerName, err))
		return false, err
	}

	if containerInspect.State.Running {
		containerLogger.WriteToFileOK(fmt.Sprintf("Container %s is running", containerName))
		return true, nil
	}
	containerLogger.WriteToFileError(fmt.Sprintf("Container %s is not running", containerName))
	return false, nil
}

func TarDirectory(dirPath string) (io.ReadCloser, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	err := filepath.Walk(dirPath, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		hdr, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}
		hdr.Name, _ = filepath.Rel(dirPath, file)
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if !fi.IsDir() {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	tw.Close()
	return io.NopCloser(buf), nil
}
