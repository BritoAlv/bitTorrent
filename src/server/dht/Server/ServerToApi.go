package main

import (
	common "Centralized/Common"
	"fmt"
	"net"
	"net/http"
	"time"
)

func handleIPRequest(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	name := r.URL.Query().Get("name")
	logger.WriteToFileOK(fmt.Sprintf("Received IP Query asking for name : %s", name))
	var response string
	if ip, exist := clientLocations[name]; exist {
		response = ip
		logger.WriteToFileOK(fmt.Sprintf("Response IP Query : Name %s has ip %s", name, ip))
	} else {
		response = common.MsgNotExists
		logger.WriteToFileOK(fmt.Sprintf("Response IP Query : Name %s does not exist", name))
	}
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(response))
	if err != nil {
		panic(err)
	}
}

func handleLoginRequest(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	name := r.URL.Query().Get("name")
	ip := r.URL.Query().Get("ip")
	logger.WriteToFileOK(fmt.Sprintf("Received Login Query with name : %s and IP = %s", name, ip))
	var response string
	if _, exist := clientLocations[name]; exist {
		response = common.MsgNameExists
		logger.WriteToFileOK(fmt.Sprintf("Response Login Query : Name %s already exists", name))
	} else {
		response = common.MsgLogged
		logger.WriteToFileOK(fmt.Sprintf("Response Login Query : Name %s has been added with IP = %s", name, ip))
	}
	clientLocations[name] = ip
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(response))
	if err != nil {
		panic(err)
	}
}

func pingClients() {
	for {
		time.Sleep(5 * time.Second)
		mu.Lock()
		for client, location := range clientLocations {
			logger.WriteToFileOK(fmt.Sprintf("Pinging client %s at location %s", client, location))
			_, err := net.Dial("tcp", location)
			if err != nil {
				logger.WriteToFileError(fmt.Sprintf("Error pinging client %s at location %s", client, location))
				logger.WriteToFileOK(fmt.Sprintf("Decision is Removing client %s at location %s", client, location))
				delete(clientLocations, client)
				break
			}
			logger.WriteToFileOK(fmt.Sprintf("Client %s at location %s, apparently, is still  alive", client, location))
		}
		mu.Unlock()
	}
}
