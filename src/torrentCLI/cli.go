package main

import (
	"bittorrent/common"
	"fmt"
)

func main() {
	fmt.Println("Basic Torrent Utility CLI")
	var filename string
	var torrentname string
	var trackerLocation string
	fmt.Print("Enter the filename: ")
	fmt.Scanln(&filename)
	fmt.Print("Enter the torrentname: ")
	fmt.Scanln(&torrentname)
	fmt.Print("Enter the tracker location: ")
	fmt.Scanln(&trackerLocation)
	err := common.CreateTorrentFile(filename, torrentname, trackerLocation, false)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Torrent file created successfully")
	}
}
