package main

import (
	"bittorrent/server/TrackerNode"
	"fmt"
	"sync"
)

func main() {
	var ip, port string
	fmt.Print("Enter IP: ")
	fmt.Scanln(&ip)
	fmt.Print("Enter port: ")
	fmt.Scanln(&port)
	// Create a tracker somewhere.
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
