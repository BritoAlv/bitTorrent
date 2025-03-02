package Scenarios

import (
	"bittorrent/common"
	"bittorrent/dht/library/InMemory"
	"bittorrent/dht/library/Manager"
	"math/rand/v2"
	"time"
)

func ScenarioMedium() Manager.IManagerRPC {
	common.SetLogDirectoryPath("MediumScenario")
	var database = *InMemory.NewDataBaseInMemory()
	go func() {
		for {
			time.Sleep(3 * time.Second)
			if rand.Float32() <= 0.4 {
				database.CreateRandomNode()
			}
			if rand.Float32() <= 0.1 && len(database.GetNodes()) > 0 {
				database.RemoveRandomNode()
			}
			if rand.Float32() <= 0.7 {
				database.PutRandomDataRandomNode()
			}
		}
	}()
	return &database
}
