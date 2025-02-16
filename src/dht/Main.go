package main

import (
	"bittorrent/dht/library"
)

func main() {
	var database = *library.NewDataBaseInMemory()
	var zero [library.NumberBits]uint8
	var server1 = library.ServerInMemory{
		DataBase:             &database,
		ServerId:             "Alvaro",
		ChannelCommunication: nil,
		NodeId:               zero,
	}

	var client1 = library.ClientInMemory{
		DataBase: &database,
	}
	database.AddServer(&server1)

	node1 := library.NewBruteChord[library.InMemoryContact](&server1, &client1, library.NewMonitorHand[library.InMemoryContact]())
	node1.BeginWorking()
}
