package main

import (
	"bittorrent/server/TrackerNode"
	"flag"
	"fmt"
	"sync"
)

func main() {
	var ip, port string
	flag.StringVar(&ip, "ip", "127.0.0.1", "IP address to bind the tracker")
	flag.StringVar(&port, "port", "8080", "Port to bind the tracker")
	flag.Parse()
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
