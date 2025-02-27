package Scenarios

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord/Core"
	"bittorrent/dht/library/InMemory"
	"bittorrent/dht/library/Manager"
	"fmt"
	"time"
)

func ScenarioPut() Manager.IManagerRPC[InMemory.ContactInMemory] {
	N := 10
	common.SetLogDirectoryPath("PutScenario")
	var database = *InMemory.NewDataBaseInMemory()
	var toPut = make(map[Core.ChordHash][]byte)
	fmt.Printf("Going to Add N = %v  Nodes", N)
	for i := 0; i < N; i++ {
		database.CreateRandomNode()
	}
	for i := 0; i < 50; i++ {
		toPut[Core.ChordHash(i)] = []byte{byte(i)}
	}
	time.Sleep(2 * time.Second)
	nodes := database.GetNodes()
	println(len(nodes))
	for key, value := range toPut {
		for _, node := range nodes {
			go node.Put(key, value)
		}
	}
	return &database
}
