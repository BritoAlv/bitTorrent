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
	mockClient1 := peer.PeerMock{
		Address: common.Address{Ip: "localhost", Port: "8085"},
	}
	mockClient2 := peer.PeerMock{
		Address: common.Address{Ip: "localhost", Port: "8090"},
	}

	// wg := sync.WaitGroup{}
	// wg.Add(1)
	// go mockClient.Torrent(&wg)
	// wg.Wait()

	go mockClient1.Torrent(nil)
	// go mockClient2.Torrent(nil)

	var cTracker tracker.Tracker = tracker.CentralizedTracker{
		Address: common.Address{
			Ip:   "localhost",
			Port: "8000",
		}}

	torrentInfo := peer.TorrentInfo{
		Hash:        []byte{},
		Name:        "",
		PieceLength: 0,
		Pieces:      []byte{},
		Length:      0,
		Files:       []peer.FileInfo{},
	}

	notificationChannel := make(chan interface{}, 1000)

	client := peer.Peer{
		Address: common.Address{
			Ip:   "localhost",
			Port: "9000",
		},
		TorrentInfo:         torrentInfo,
		Tracker:             cTracker,
		NotificationChannel: notificationChannel,
		Peers:               make(map[common.Address]peer.PeerInfo),
	}

	peerGroup := sync.WaitGroup{}
	peerGroup.Add(1)
	go client.Torrent(&peerGroup)

	time.Sleep(time.Second * time.Duration(10))
	go mockClient2.ActiveTorrent(nil, client.Address)
	// notificationChannel <- peer.KillNotification{}

	peerGroup.Wait()
	fmt.Println(client.Peers)
}

// func Connect(torrentFilePath string) {
// 	// Decode the torrent file
// 	// Verify info-hash
// 	// Create a tracker depending on the specification of .torrent (centralized or not)

// 	// Build track-request
// 	// Send track-request => tracker.Track(tracker.TrackRequest{})
// 	// Grab peers

// 	//** Concurrent things to do
// 	// 1. Send a request to tracker every an specified time
// 	// 2. If the number of peers drops below 20 send a request to tracker
// 	//
// }
