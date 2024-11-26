package main

import (
	"bittorrent/common"
	"bittorrent/peer"
	"bittorrent/tracker"
	"fmt"
	"sync"
	"time"
)

func main() {
	torrent, err := common.ParseTorrentFile("server.py.torrent")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	mockClient1 := peer.PeerMock{
		Address: common.Address{Ip: "localhost", Port: "8085"},
	}
	mockClient2 := peer.PeerMock{
		Address: common.Address{Ip: "localhost", Port: "8090"},
	}

	go mockClient1.Torrent(nil)
	go mockClient2.Torrent(nil)

	address := common.Address{Ip: "localhost", Port: "9000"}
	var centralizedTracker tracker.Tracker = tracker.CentralizedTracker{
		Url: torrent.Announce,
	}

	client := peer.New(
		"peer-id",
		address,
		torrent,
		centralizedTracker,
	)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go client.Torrent(&waitGroup)

	time.Sleep(time.Second * time.Duration(10))
	go mockClient2.ActiveTorrent(nil, address)

	waitGroup.Wait()
}
