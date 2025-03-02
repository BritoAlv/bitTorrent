package main

import (
	"bittorrent/dht/library/Manager"
	"bittorrent/dht/library/WithSocket"
	"strconv"
)

func main() {
	WithSocket.RegisterStartUp("tun0", "GUI", []string{"12345", "12346", "12347", "12348", "12349"})
	portsExposed := make([]string, 0)
	for i := 0; i < 20; i++ {
		portsExposed = append(portsExposed, strconv.Itoa(9200+i))
	}
	manager := Manager.NewHttpManager(portsExposed)
	gui := Manager.NewGUI(manager)
	gui.Start()
}
