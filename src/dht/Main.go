package main

import (
	"bittorrent/dht/library/InMemory"
	"bittorrent/dht/library/Manager"
	"bittorrent/dht/library/Scenarios"
)

func main() {
	manager := Scenarios.ScenarioEasy()
	gui := Manager.NewGUI[InMemory.ContactInMemory](&manager)
	gui.Start()
}
