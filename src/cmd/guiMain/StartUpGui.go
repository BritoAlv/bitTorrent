package main

import (
	"bittorrent/dht/library/Manager"
	"bittorrent/dht/library/WithSocket"
)

func main() {
	WithSocket.RegisterStartUp("tun0", "GUI", []string{"12345", "12346", "12347", "12348", "12349"})
	manager := WithSocket.NewManagerSocket()
	gui := Manager.NewGUI[WithSocket.SocketContact](manager)
	gui.Start()
}
