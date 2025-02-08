package main

import (
	"bittorrent/common"
	"fmt"
	"os"
	"os/exec"
)

var logger = common.NewLogger("./build.txt")

func buildProject(dir string, output string) {
	logger.WriteToFileOK(fmt.Sprintf("Building the project in directory %s with output %s", dir, output))
	cmd := exec.Command("go", "build", "-o", output)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		logger.WriteToFileError("Error building project: " + err.Error())
	}
	logger.WriteToFileOK(fmt.Sprintf("Done building the project"))
}

func main() {
	buildProject("client", "client")
	buildProject("server", "server")
	buildProject("torrentCLI", "cli")
}
