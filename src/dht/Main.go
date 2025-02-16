package main

import (
	"bittorrent/dht/library"
)

func main() {
	var database = *library.NewDataBaseInMemory()
	var server1 = library.NewServerInMemory(&database, "BritoServer")

	var client1 = library.NewClientInMemory(&database)
	database.AddServer(server1)

	node1 := library.NewBruteChord[library.InMemoryContact](server1, client1, library.NewMonitorHand[library.InMemoryContact]())
	node1.BeginWorking()
}
