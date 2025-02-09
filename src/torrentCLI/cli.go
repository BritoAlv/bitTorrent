package main

import (
	"bittorrent/torrent"
	"fmt"
)

func main() {
	fmt.Println("Basic Torrent Utility CLI")
	var fileName string
	var torrentName string
	var trackerLocation string
	fmt.Print("Enter the fileName: ")
	fmt.Scanln(&fileName)
	fmt.Print("Enter the torrentName: ")
	fmt.Scanln(&torrentName)
	fmt.Print("Enter the tracker location: ")
	fmt.Scanln(&trackerLocation)
	err := torrent.CreateTorrentFile(fileName, torrentName, trackerLocation)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Torrent file created successfully")
	}
}
