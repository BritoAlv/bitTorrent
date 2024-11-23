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
	torrent, err := common.ParseTorrentFile("ubuntu-22.04.5-desktop-amd64.iso.torrent")

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

	var centralizedTracker tracker.Tracker = tracker.CentralizedTracker{
		Url: "http://localhost:5000/tracker",
	}

	notificationChannel := make(chan interface{}, 1000)

	client := peer.Peer{
		Id: "peer_id",
		Address: common.Address{
			Ip:   "localhost",
			Port: "9000",
		},
		TorrentData:         torrent,
		Tracker:             centralizedTracker,
		NotificationChannel: notificationChannel,
		Peers:               make(map[common.Address]peer.PeerInfo),
	}

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go client.Torrent(&waitGroup)

	time.Sleep(time.Second * time.Duration(10))
	go mockClient2.ActiveTorrent(nil, client.Address)
	// notificationChannel <- peer.KillNotification{}

	waitGroup.Wait()
	fmt.Println(client.Peers)
}
