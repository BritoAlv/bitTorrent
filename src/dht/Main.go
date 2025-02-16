package main

import (
	"bittorrent/dht/library"
	"strconv"
	"sync"
)

func main() {
	var database = *library.NewDataBaseInMemory()
	var barrier = sync.WaitGroup{}
	N := 24
	for i := 0; i < N; i++ {
		iString := strconv.Itoa(i)
		var server = library.NewServerInMemory(&database, "Server"+iString)
		var client = library.NewClientInMemory(&database, "Client"+iString)
		database.AddServer(server)
		node := library.NewBruteChord[library.InMemoryContact](server, client, library.NewMonitorHand[library.InMemoryContact]("Monitor"+iString))
		go func() {
			barrier.Add(1)
			node.BeginWorking()
			barrier.Done()
		}()
	}
	barrier.Wait()
}
