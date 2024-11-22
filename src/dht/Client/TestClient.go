package main

import (
	"Centralized/Common"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"sync"
)

const SizeBinaryString = 10
const ServerUrl = "http://localhost:8080/"
const requests = 10
const stringsPerRequest = 4

/*
Ensure server is consistent :
	If client ask for a name, then it should be added or reported that exist,
	but once an answer of this type is obtained, if this name gets asked again
	then server can-t say it's new.
*/

var stored = map[string]bool{}
var lock = sync.Mutex{}
var group = sync.WaitGroup{}

func sendGetRequest(direction string, argStr string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(ServerUrl+direction+"?arg=%s", argStr))
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

func RequestSender(logger *common.Logger) {
	strings := [stringsPerRequest]string{}
	confirmed := [stringsPerRequest]bool{}
	isNew := [stringsPerRequest]bool{}
	lock.Lock()
	for i := 0; i < stringsPerRequest; i++ {
		str, err := common.GenerateRandomString(SizeBinaryString)
		if err != nil {
			logger.WriteToFileError("Error when generating random string : " + err.Error())
			return
		}
		strings[i] = str
		confirmed[i] = stored[str]
		logger.WriteToFileOK(fmt.Sprintf("Request generated string %s and IsConfirmed == %t ", str, confirmed[i]))
	}
	lock.Unlock()
	for i := 0; i < stringsPerRequest; i++ {
		if rand.Float32() <= 0.5 {
			logger.WriteToFileOK(fmt.Sprintf("Request %d Login string %s  ", i, strings[i]))
			result, err := sendGetRequest(common.LoginURL, strings[i])
			logger.WriteToFileOK(fmt.Sprintf("Request %d Login string %s  received answer %s ", i, strings[i], result))
			if err != nil {
				logger.WriteToFileError(err.Error())
				return
			}
			isNew[i] = result == common.MsgLogged
			logger.WriteToFileOK(fmt.Sprintf("string %s isNew == %t ", strings[i], isNew[i]))
		} else {
			logger.WriteToFileOK(fmt.Sprintf("Request %d IP  of string %s  ", i, strings[i]))
			result, err := sendGetRequest(common.IPURL, strings[i])
			logger.WriteToFileOK(fmt.Sprintf("Request %d IP string %s  received answer %s ", i, strings[i], result))
			if err != nil {
				logger.WriteToFileError(err.Error())
				return
			}
			if result == common.MsgNotExists {
				isNew[i] = true
			}
			logger.WriteToFileOK(fmt.Sprintf("string %s isNew == %t ", strings[i], isNew[i]))
		}
	}
	for i := 0; i < stringsPerRequest; i++ {
		if confirmed[i] && isNew[i] {
			logger.WriteToFileError(fmt.Sprintf("string %s was previously added but server says its new ", strings[i]))
			return
		} else if !confirmed[i] {
			lock.Lock()
			logger.WriteToFileOK(fmt.Sprintf("string %s is confirmed ", strings[i]))
			stored[strings[i]] = true
			lock.Unlock()
		}
	}
}

func TestClient() {
	for i := 0; i < requests; i++ {
		group.Add(1)
		go func() {
			var logger = common.NewLogger(fmt.Sprintf("ClientLog%d.txt", i+1))
			defer group.Done()
			RequestSender(logger)
			logger.WriteToFileOK(fmt.Sprintf("Server Requests %d finished !!! ", i+1))
		}()
	}
	group.Wait()
}
