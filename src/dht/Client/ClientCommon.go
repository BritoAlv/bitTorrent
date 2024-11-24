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

func tryCommunicate(local string, friend string, logger *common.Logger) bool {
	logger.WriteToFileOK(fmt.Sprintf("Requesting IP of %s friend", friend))
	response, err := sendGetRequest(common.BuildIPRequest(ServerUrl, friend))
	if err != nil {
		logger.WriteToFileError(err.Error())
		return false
	}
	if response == common.MsgNotExists {
		logger.WriteToFileOK(fmt.Sprintf("Server doesn't know about %s ", friend))
		return false
	}
	logger.WriteToFileOK(fmt.Sprintf("Obtained from server IP of %s : %s", friend, response))
	friendAddress := response
	conn, err := net.Dial("tcp", friendAddress)
	if err != nil {
		logger.WriteToFileError(err.Error())
		return false
	}
	logger.WriteToFileOK(fmt.Sprintf("Connected to %s", friend))
	_, err = conn.Write([]byte(fmt.Sprintf("Hey I'm %s and want to communicate with you", local)))
	if err != nil {
		logger.WriteToFileError("Error sending hello message: " + err.Error())
		return false
	}
	logger.WriteToFileOK(fmt.Sprintf("Sent hello message to %s", friend))
	return true
}
