package main

import (
	"bittorrent/common"
	"bittorrent/server/TrackerNode"
	"fmt"
	"sync"
)

func main() {
	// Create a tracker somewhere.
	tracker1 := TrackerNode.NewHttpTracker("1234", "MyHTTPServer")
	// Create a torrent file with the tracker's URL.
	err := common.CreateTorrentFile("main.go", "main1", tracker1.SaveTorrent(), false)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Start the tracker. After this peers may be able to find the tracker.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = tracker1.Listen()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	wg.Wait()
}