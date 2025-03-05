package InMemory

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord/Core"
	"bytes"
	"math/rand/v2"
	"testing"
	"time"
)

var expectedHashTable = map[Core.ChordHash][]byte{}

func putData(t *testing.T, database *DataBaseInMemory) {
	value := []byte(common.GenerateRandomString(10))
	key := Core.GenerateRandomBinaryId()
	t.Logf("Generated key %v with value %v to add to the HashTable", key, value)
	randomNode := database.getRandomNode(1)[0]
	t.Logf("Random node with ID = %v will be used to call do the Put RPC", randomNode.GetId())
	result := randomNode.Put(key, value)
	if result {
		expectedHashTable[key] = value
		t.Logf("Added the value to the hash table and to Chord")
	}
}

func updateData(t *testing.T, database *DataBaseInMemory) {
	for key, value := range expectedHashTable {
		valueUpdated := []byte(common.GenerateRandomString(10))
		t.Logf("Going to update key %v with value %v to new value %v ", key, value, valueUpdated)
		randomNode := database.getRandomNode(1)[0]
		t.Logf("Random node with ID = %v will be used to call do the Put RPC", randomNode.GetId())
		result := randomNode.Put(key, valueUpdated)
		if result {
			expectedHashTable[key] = valueUpdated
			t.Logf("Added the value to the hash table and to Chord")
		}
		break
	}
}

func getCheck(t *testing.T, database *DataBaseInMemory) {
	for key, value := range expectedHashTable {
		t.Logf("Going to get from Chord value of key %v", key)
		randomNode := database.getRandomNode(1)[0]
		t.Logf("Random node with ID = %v will be used to call do the Get RPC", randomNode.GetId())
		valueFromChord, exist := randomNode.Get(key)
		if !exist {
			t.Fatalf("Could not find the key %v in Chord, but exist", key)
		}
		if !bytes.Equal(valueFromChord, value) {
			t.Fatalf("Chord value %v does not match the expected value %v in Chord", valueFromChord, value)
		}
	}
}

func removeTwoNodes(t *testing.T, database *DataBaseInMemory) {
	if len(database.GetNodes()) >= 3 {
		ids := database.RemoveRandomNode(2)
		t.Logf("Killed two nodes with ID = %v, %v", ids[0], ids[1])
	}
}

func removeOneNodes(t *testing.T, database *DataBaseInMemory) {
	if len(database.GetNodes()) >= 2 {
		ids := database.RemoveRandomNode(1)
		t.Logf("Killed one node with ID = %v", ids[0])
	}
}

func addOneNode(t *testing.T, database *DataBaseInMemory) {
	if len(database.GetNodes()) <= 2 {
		id := database.CreateRandomNode()
		t.Logf("Added one node with ID = %v", id)
	}
}

func TestWithSimulation(t *testing.T) {
	database, _, _ := StartUp("TestWithSimulation", 10)
	expectedHashTable = make(map[Core.ChordHash][]byte)
	events := []func(*testing.T, *DataBaseInMemory){
		putData,
		updateData,
		getCheck,
		removeTwoNodes,
		removeOneNodes,
		addOneNode,
	}
	timeout := time.After(60 * time.Minute)
	for {
		select {
		case <-timeout:
			t.Logf("TestWithSimulation passed")
			return

		default:
			time.Sleep(5 * time.Second)
			randEvent := rand.Int() % len(events)
			events[randEvent](t, database)
		}
	}
}
