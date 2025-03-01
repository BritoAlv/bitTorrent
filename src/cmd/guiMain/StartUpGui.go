package main

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord/Core"
	"bittorrent/dht/library/Manager"
	"bittorrent/dht/library/WithSocket"
)

func main() {
	Core.RegisterNotifications[WithSocket.SocketContact]()
	common.SetLogDirectoryPath("./GUI")
	manager := WithSocket.NewManagerSocket()
	gui := Manager.NewGUI[WithSocket.SocketContact](manager)
	gui.Start()
}
