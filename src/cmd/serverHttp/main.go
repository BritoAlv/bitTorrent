package main

import (
	"bittorrent/dht/library/WithSocket"
	"bittorrent/server/TrackerNode"
	"fmt"
	"time"
)

func main() {
	iface := "eth0"
	WithSocket.RegisterStartUp(iface, "HttpChord", []string{"12345"})
	var tracker1 = TrackerNode.NewHttpTracker("TrackerDocker", iface, "8080")
	for {
		time.Sleep(1 * time.Second)
		fmt.Printf("Tracker is running %v:%v \n", tracker1.Ip, tracker1.Port)
	}
}
