package main

import (
	"Centralized/Common"
	"fmt"
	"net/http"
	"sync"
)

var (
	names = map[string]bool{
		common.RootUser: true,
	}
	mu sync.Mutex
)

var queryID = 0
var logger = common.NewLogger("ServerLog.txt")

func main() {
	logger.WriteToFile("Server started")
	http.HandleFunc("/", handleRequest)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	queryID++
	arg := r.URL.Query().Get("arg")
	logger.WriteToFile(fmt.Sprintf("Received query %d with name : %s", queryID, arg))
	var response string
	if _, exist := names[arg]; exist {
		response = common.MsgNameExists
		logger.WriteToFile(fmt.Sprintf("Response to query %d : Name %s already exists", queryID, arg))
		return
	} else {
		names[arg] = true
		response = common.MsgLogged
		logger.WriteToFile(fmt.Sprintf("Response to query %d : Name %s added", queryID, arg))
	}
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(response))
	if err != nil {
		panic(err)
	}
}
