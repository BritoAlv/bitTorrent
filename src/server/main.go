package main

import (
	"bittorrent/dht/library/WithSocket"
	"bittorrent/server/TrackerNode"
	"fmt"
	"sync"
)

func main() {

	ip, _ := WithSocket.GetIpFromInterface("eth0")
	if ip == "" {
		ip = "localhost"
	}
	port := "8080"

	var tracker1 = TrackerNode.NewHttpTracker(ip+":"+port, "MyHTTPServer")
	fmt.Printf("Tracker Location is : %s\n", tracker1.SaveTorrent())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := tracker1.Listen()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	wg.Wait()
}
