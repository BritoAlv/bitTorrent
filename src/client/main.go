package main

import (
	"bittorrent/client/peer"
	"bittorrent/common"
	"fmt"
	"sync"
)

func main() {
	var ip, port string
	fmt.Print("Enter IP: ")
	fmt.Scanln(&ip)
	fmt.Print("Enter port: ")
	fmt.Scanln(&port)
	torrent, err := common.ParseTorrentFile("./main1.torrent")
	if err != nil{
		fmt.Println(err)
		return
	}

	peer, err := peer.New(common.GenerateRandomString(20), common.Address{
		Ip:   ip,
		Port: port,
	}, torrent, "./")

	if err != nil{
		fmt.Println(err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	peer.Torrent(&wg)
	wg.Wait()
}