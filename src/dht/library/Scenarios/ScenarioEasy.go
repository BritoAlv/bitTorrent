package Scenarios

import (
	"bittorrent/common"
	"bittorrent/dht/library/InMemory"
	"bittorrent/dht/library/Manager"
	"fmt"
	"math/rand/v2"
	"time"
)

func ScenarioEasy() Manager.IManagerRPC {
	common.SetLogDirectoryPath("BasicScenario")
	var database = *InMemory.NewDataBaseInMemory()
	fmt.Println("Nodes are being added and removed randomly every once a while")
	go func() {
		for {
			time.Sleep(1 * time.Second)
			if rand.Float32() <= 0.3 {
				database.CreateRandomNode()
			}
			if rand.Float32() <= 0.1 && len(database.GetNodes()) > 0 {
				database.RemoveRandomNode()
			}
		}
	}()
	return &database
}
