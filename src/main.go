package main

import (
	"bittorrent/common"
	"bittorrent/peer"
	"fmt"
	"sync"
)

func main() {
	torrent, err := common.ParseTorrentFile("regex.pdf.torrent")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// mockClient1 := peer.PeerMock{
	// 	Address: common.Address{Ip: "localhost", Port: "8085"},
	// }
	// mockClient2 := peer.PeerMock{
	// 	Address: common.Address{Ip: "localhost", Port: "8090"},
	// }

	// go mockClient1.Torrent(nil)
	// go mockClient2.Torrent(nil)

	id := "peer-id"
	address := common.Address{Ip: "localhost", Port: "9000"}

	client, err := peer.New(id, address, torrent)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go client.Torrent(&waitGroup)

	// time.Sleep(time.Second * time.Duration(10))
	// go mockClient2.ActiveTorrent(nil, address)

	waitGroup.Wait()
}
