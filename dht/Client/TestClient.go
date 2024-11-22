package main

import (
	"Centralized/Common"
	"fmt"
	"io"
	"net/http"
	"sync"
)

const SizeBinaryString = 160
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

func sendRequest(str string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(ServerUrl+"?arg=%s", str))
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

func RequestsSender(id int) {
	var logger = common.NewLogger(fmt.Sprintf("ClientLog%d.txt", id))
	strings := [stringsPerRequest]string{}
	confirmed := [stringsPerRequest]bool{}
	isNew := [stringsPerRequest]bool{}
	lock.Lock()
	for i := 0; i < stringsPerRequest; i++ {
		str, err := common.GenerateRandomString(SizeBinaryString)
		if err != nil {
			logger.WriteToFile("Error when generating random string : " + err.Error())
			return
		}
		strings[i] = str
		confirmed[i] = stored[str]
		logger.WriteToFile(fmt.Sprintf("Request generated string %s and IsConfirmed == %t ", str, confirmed[i]))
	}
	lock.Unlock()
	for i := 0; i < stringsPerRequest; i++ {
		logger.WriteToFile(fmt.Sprintf("Request sending string %s %d ", strings[i], i))
		result, err := sendRequest(strings[i])
		logger.WriteToFile(fmt.Sprintf("Request from string %s %d received answer %s ", strings[i], i, result))
		if err != nil {
			logger.WriteToFile(err.Error())
			return
		}
		isNew[i] = result == common.MsgLogged
		logger.WriteToFile(fmt.Sprintf("string %s isNew == %t ", strings[i], isNew[i]))
	}
	for i := 0; i < stringsPerRequest; i++ {
		if confirmed[i] && isNew[i] {
			logger.WriteToFile(fmt.Sprintf("string %s was previously added but server says its new ", strings[i]))
			return
		} else if !confirmed[i] {
			lock.Lock()
			logger.WriteToFile(fmt.Sprintf("string %s is confirmed ", strings[i]))
			stored[strings[i]] = true
			lock.Unlock()
		}
	}
	logger.WriteToFile(fmt.Sprintf("Request %d finished successfully !!! ", id))
}

func main() {
	for i := 0; i < requests; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			RequestsSender(i + 1)
		}()
	}
	group.Wait()
}
