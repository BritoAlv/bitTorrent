package main

import (
	"bittorrent/dht/library/Manager"
	"strconv"
)

func main() {
	portsExposed := make([]string, 0)
	for i := 0; i < 20; i++ {
		portsExposed = append(portsExposed, strconv.Itoa(9200+i))
	}
	manager := Manager.NewHttpManager(portsExposed)
	gui := Manager.NewGUI(manager)
	gui.Start()
}
