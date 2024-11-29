package main

import (
	"bittorrent/common"
	"fmt"
	"net/http"
)

func main() {
	trackerUrl := "http://127.0.0.1:1234/announce"
	err := common.CreateTorrentFile("main.go", trackerUrl, false)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	torrent, err := common.ParseTorrentFile("./torrents/main.go.torrent")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var request common.TrackRequest
	request.PeerId = "Alvaro"
	request.InfoHash = torrent.InfoHash
	request.Ip = "127.0.0.5"
	request.Port = "332"
	request.Left = int(torrent.Length)
	encoded, err := common.BuildHttpUrl(trackerUrl, request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	http.Get(encoded)
}
