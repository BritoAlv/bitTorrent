package main

import (
	"Centralized/Common"
	"fmt"
	"net/http"
	"sync"
)

var (
	clientLocations = map[string]string{}
	mu              sync.Mutex
)

var queryID = 0
var logger = common.NewLogger("ServerLog.txt")

func main() {
	logger.WriteToFileOK("Server started")
	http.HandleFunc("/"+common.LoginURL, handleLoginRequest)
	http.HandleFunc("/"+common.IPURL, handleIPRequest)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func handleIPRequest(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	queryID++
	arg := r.URL.Query().Get("arg")
	logger.WriteToFileOK(fmt.Sprintf("Received IP Query %d asking for name : %s", queryID, arg))
	var response string
	if ip, exist := clientLocations[arg]; exist {
		response = ip
		logger.WriteToFileOK(fmt.Sprintf("Response IP Query %d : Name %s has ip %s", queryID, arg, ip))
	} else {
		response = common.MsgNotExists
		logger.WriteToFileOK(fmt.Sprintf("Response IP Query %d : Name %s does not exist", queryID, arg))
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
	arg := r.URL.Query().Get("arg")
	logger.WriteToFileOK(fmt.Sprintf("Received Login Query %d with name : %s", queryID, arg))
	var response string
	clientLocations[arg] = r.RemoteAddr
	if _, exist := clientLocations[arg]; exist {
		response = common.MsgNameExists
		logger.WriteToFileOK(fmt.Sprintf("Response Login Query %d : Name %s already exists", queryID, arg))
		return
	} else {
		response = common.MsgLogged
		logger.WriteToFileOK(fmt.Sprintf("Response Login Query %d : Name %s added", queryID, arg))
	}
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(response))
	if err != nil {
		panic(err)
	}
}
