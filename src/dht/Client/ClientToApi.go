package main

import (
	common "Centralized/Common"
	"fmt"
	"net"
	"time"
)

const numberOfFriends = 5

/*
Each client will be listening for upcoming connections from friends.
*/
func listen(listener *net.Listener, name string, logger *common.Logger) {
	for {
		conn, err := (*listener).Accept()
		if err != nil {
			logger.WriteToFileError("Accepting connection failed: " + err.Error())
			continue
		}
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			logger.WriteToFileError("Reading from connection Failed: " + err.Error())
			continue
		}
		logger.WriteToFileOK(fmt.Sprintf("Received message from %s: %s", name, string(buf[:n])))
	}
}

/*
Each client will have someFriends whom he wishes to send a message.
*/
func speak(source string, logger *common.Logger) {
	friends := [numberOfFriends]string{}
	communicated := [numberOfFriends]bool{}
	for i := 0; i < numberOfFriends; i++ {
		str, err := common.GenerateRandomString(SizeBinaryString)
		if err != nil {
			logger.WriteToFileError("When generating random string: " + err.Error())
			return
		}
		friends[i] = str
		logger.WriteToFileOK(fmt.Sprintf("Save %s as a candidate client to communicate with", str))
	}
	for {
		time.Sleep(2 * time.Second)
		var flag = true
		for i := 0; i < numberOfFriends; i++ {
			flag = flag && communicated[i]
		}
		if flag {
			logger.WriteToFileOK(fmt.Sprintf("Already done with communications congrats !!!"))
			break
		}
		for i := 0; i < numberOfFriends; i++ {
			if !communicated[i] {
				communicated[i] = tryCommunicate(source, friends[i], logger)
			}
		}
	}
}

/*
when client arrives it has to send a request to the server to sign-up and tell to the others
its location / ipaddress.
*/
func register(address string, logger *common.Logger) (string, error) {
	for {
		name, err := common.GenerateRandomString(SizeBinaryString)
		if err != nil {
			logger.WriteToFileError("When generating random string: " + err.Error())
			return "", err
		}
		logger.WriteToFileOK(fmt.Sprintf("Trying to login with name %s", name))
		result, err := sendGetRequest(common.BuildLoginRequest(ServerUrl, name, address))
		if err != nil {
			logger.WriteToFileError(err.Error())
			return "", err
		}
		if result == common.MsgNameExists {
			logger.WriteToFileOK(fmt.Sprintf("Name %s already exists", name))
			name, err = common.GenerateRandomString(SizeBinaryString)
			if err != nil {
				logger.WriteToFileError("When generating random string: " + err.Error())
			}
		} else if result == common.MsgLogged {
			logger.WriteToFileOK(fmt.Sprintf("Name %s registered", name))
			return name, nil
		} else {
			logger.WriteToFileError(fmt.Sprintf("Unexpected response from server: %s", result))
			return "", fmt.Errorf("unexpected response from server: %s", result)
		}
	}
}