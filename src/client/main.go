package main

import (
	"bittorrent/client/peer"
	"bittorrent/common"
	"fmt"
	"sync"
)

func main() {
	fmt.Println("Given an IP, port and torrentPath this will start a peer bounded to those")
	var ip, port, torrentPath string
	fmt.Print("Enter IP: ")
	fmt.Scanln(&ip)
	fmt.Print("Enter port: ")
	fmt.Scanln(&port)
	fmt.Print("Enter torrentPath: ")
	fmt.Scanln(&torrentPath)

	torrent, err := common.ParseTorrentFile(torrentPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	peer, err := peer.New(common.GenerateRandomString(20), common.Address{
		Ip:   ip,
		Port: port,
	}, torrent, "./")

	if err != nil {
		fmt.Println(err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	peer.Torrent(&wg)
	wg.Wait()
}
