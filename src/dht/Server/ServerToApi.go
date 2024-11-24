package main

import (
	common "Centralized/Common"
	"fmt"
	"net/http"
)

var queryID = 0 // Used to identify the queries for logging purposes.

func handleIPRequest(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	queryID++
	name := r.URL.Query().Get("name")
	logger.WriteToFileOK(fmt.Sprintf("Received IP Query %d asking for name : %s", queryID, name))
	var response string
	if ip, exist := clientLocations[name]; exist {
		response = ip
		logger.WriteToFileOK(fmt.Sprintf("Response IP Query %d : Name %s has ip %s", queryID, name, ip))
	} else {
		response = common.MsgNotExists
		logger.WriteToFileOK(fmt.Sprintf("Response IP Query %d : Name %s does not exist", queryID, name))
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
	queryID++
	name := r.URL.Query().Get("name")
	ip := r.URL.Query().Get("ip")
	logger.WriteToFileOK(fmt.Sprintf("Received Login Query %d with name : %s and IP = %s", queryID, name, ip))
	var response string
	if _, exist := clientLocations[name]; exist {
		response = common.MsgNameExists
		logger.WriteToFileOK(fmt.Sprintf("Response Login Query %d : Name %s already exists", queryID, name))
	} else {
		response = common.MsgLogged
		logger.WriteToFileOK(fmt.Sprintf("Response Login Query %d : Name %s has been added with IP = %s", queryID, name, ip))
	}
	clientLocations[name] = ip
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(response))
	if err != nil {
		panic(err)
	}
}
