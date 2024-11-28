package main

import (
	"Centralized/Common"
	"net/http"
	"sync"
)

var (
	clientLocations = map[string]string{}
	mu              sync.Mutex
)

var logger = common.NewLogger("ServerLog.txt")

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		logger.WriteToFileOK("Server started")
		http.HandleFunc("/"+common.LoginRoute, handleLoginRequest)
		http.HandleFunc("/"+common.IPRoute, handleIPRequest)
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()
	go func() { pingClients() }()
	wg.Wait()
}
