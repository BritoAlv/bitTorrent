package main

import (
	common "Centralized/Common"
	"fmt"
	"io"
	"net"
	"net/http"
)

const SizeBinaryString = 3
const ServerUrl = "http://localhost:8080/"

func sendGetRequest(body string) (string, error) {
	resp, err := http.Get(body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("resp.StatusCode was %d", resp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	answer := string(bodyBytes)
	return answer, nil
}

func tryCommunicate(source string, s string, logger *common.Logger) bool {
	logger.WriteToFileOK(fmt.Sprintf("Requesting IP of %s", s))
	response, err := sendGetRequest(common.BuildIPRequest(ServerUrl, s))
	if err != nil {
		logger.WriteToFileError(err.Error())
		return false
	}
	if response == common.MsgNotExists {
		logger.WriteToFileOK(fmt.Sprintf("Server doesn't know about %s ", s))
		return false
	}
	logger.WriteToFileOK(fmt.Sprintf("Received IP of %s: %s", s, response))
	friendAddress := response
	conn, err := net.Dial("tcp", friendAddress)
	if err != nil {
		logger.WriteToFileError(err.Error())
		return false
	}
	logger.WriteToFileOK(fmt.Sprintf("Connected to %s", s))
	_, err = conn.Write([]byte(fmt.Sprintf("Hello from %s", source)))
	if err != nil {
		logger.WriteToFileError("Error sending hello message: " + err.Error())
		return false
	}
	logger.WriteToFileOK(fmt.Sprintf("Sent hello message to %s", s))
	return true
}
