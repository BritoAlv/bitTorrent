package main

import (
	"bittorrent/dht/library/Manager"
	"bittorrent/dht/library/WithSocket"
)

func main() {
	WithSocket.RegisterStartUp("tun0", "GUI", []string{"12345", "12346", "12347", "12348", "12349"})
	manager := Manager.NewHttpManager(
		[]string{"9201", "9202", "9203", "9204", "9205", "9206", "9207", "9208"})
	gui := Manager.NewGUI(manager)
	gui.Start()
}
