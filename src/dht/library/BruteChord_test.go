package library

import (
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"
)

const NumberOfRuns = 2

func Sort(ids []ChordHash) {
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
}

func TestRunMultipleTimes(t *testing.T) {
	for i := 0; i < NumberOfRuns; i++ {
		t.Run("TestBasicChordBehaviourInitialization", TestBasicChordBehaviourInitialization)
		t.Run("TestBasicChordBehaviourNoDead", TestBasicChordBehaviourNoDead)
		t.Run("TestBasicChordBehaviourStabilization", TestBasicChordBehaviourStabilization)
		t.Run("TestBasicPutGet", TestBasicPutGet)
		t.Run("TestBasicReplication", TestBasicReplication)
		t.Run("TestReplication", TestReplication)
	}
}

func StartUp(name string, NumberNodes int) (*DataBaseInMemory, *sync.WaitGroup, map[ChordHash]*BruteChord[InMemoryContact], []ChordHash) {
	SetLogDirectoryPath(name)
	var database = *NewDataBaseInMemory()
	var barrier = sync.WaitGroup{}
	for i := 0; i < NumberNodes; i++ {
		randomId := GenerateRandomBinaryId()
		randomIdStr := strconv.Itoa(int(randomId))
		var server = NewServerInMemory(&database, "Server"+randomIdStr)
		var client = NewClientInMemory(&database, "Client"+randomIdStr)
		var monitor = NewMonitorHand[InMemoryContact]("Monitor" + randomIdStr)
		node := NewBruteChord[InMemoryContact](server, client, monitor, randomId)
		database.AddNode(node, server, client)
		barrier.Add(1)
		go func() {
			node.BeginWorking()
			barrier.Done()
		}()
	}
	nodes := make(map[ChordHash]*BruteChord[InMemoryContact])
	ids := make([]ChordHash, 0)
	for _, node := range database.GetNodes() {
		nodes[node.GetId()] = node
		ids = append(ids, node.GetId())
	}
	Sort(ids)
	time.Sleep((10 * WaitingTime) * time.Second)
	return &database, &barrier, nodes, ids
}

/*
TestBasicChordBehaviourInitialization : N nodes are created simultaneously, eventually after stabilization occurs all the
nodes should have as its successor the next one in the ring.
*/
func TestBasicChordBehaviourInitialization(t *testing.T) {
	NumberNodes := 10
	_, barrier, nodes, ids := StartUp("TestBasicChordBehaviourInitialization", NumberNodes)
	for i := 0; i < NumberNodes; i++ {
		nodeId := ids[i]
		node := nodes[nodeId]
		successor := node.GetContact(1)
		successorId := successor.NodeId
		expectedSuccessorId := ids[(i+1)%NumberNodes]
		if successorId != expectedSuccessorId {
			t.Errorf("Node %v has successor %v, expected %v", nodeId, successorId, expectedSuccessorId)
		}
	}
	for i := 0; i < NumberNodes; i++ {
		node := nodes[ids[i]]
		node.StopWorking()
	}
	barrier.Wait()
}

// TestBasicChordBehaviourNoDead : N nodes are created simultaneously, at any moment all nodes are active, so there should be no dead nodes.
func TestBasicChordBehaviourNoDead(t *testing.T) {
	_, barrier, nodes, _ := StartUp("TestBasicChordBehaviourNoDead", 10)
	time.Sleep((10 * WaitingTime) * time.Second)
	for _, node := range nodes {
		if len(node.DeadContacts) > 0 {
			t.Errorf("Node %v has the dead contacts ", node.GetId())
		}
	}
	for _, node := range nodes {
		node.StopWorking()
	}
	barrier.Wait()
}

// TestBasicChordBehaviourStabilization : N nodes are created simultaneously, some nodes randomly go down, and eventually go up, after stabilization occurs it should
// happen that all the nodes are successors are fine.
func TestBasicChordBehaviourStabilization(t *testing.T) {
	NumberNodes := 10
	_, barrier, nodes, ids := StartUp("TestBasicChordBehaviourStabilization", NumberNodes)
	down := make([]*BruteChord[InMemoryContact], 0, NumberNodes)
	for i := 0; i < NumberNodes; i++ {
		if rand.Float32() <= 0.5 {
			nodeId := ids[i]
			node := nodes[nodeId]
			down = append(down, node)
			node.StopWorking()
		}
	}
	for _, node := range down {
		barrier.Add(1)
		go func() {
			node.BeginWorking()
			barrier.Done()
		}()
	}
	time.Sleep((10 * WaitingTime) * time.Second)
	for i := 0; i < NumberNodes; i++ {
		nodeId := ids[i]
		node := nodes[nodeId]
		successor := node.GetContact(1)
		successorId := successor.NodeId
		expectedSuccessorId := ids[(i+1)%NumberNodes]
		if successorId != expectedSuccessorId {
			t.Errorf("Node %v has successor %v, expected %v", nodeId, successorId, expectedSuccessorId)
		}
	}
	for i := 0; i < NumberNodes; i++ {
		node := nodes[ids[i]]
		node.StopWorking()
	}
	barrier.Wait()
}

func TestBasicPutGet(t *testing.T) {
	NumberNodes := 10
	_, barrier, nodes, ids := StartUp("TestBasicPutGet", NumberNodes)
	// Put data into a random node
	key := GenerateRandomBinaryId()
	value := []byte("test value")
	randomIndex := rand.Intn(len(ids))
	randomNode := nodes[ids[randomIndex]]
	randomNode.Put(key, value)

	// Verify that the data is stored correctly by querying every node.
	for _, node := range nodes {
		storedValue, exist := node.Get(key)
		if exist && string(storedValue) == string(value) {
			t.Logf("Chord has the correct value for key %v", key)
		}
		if exist && !(string(storedValue) == string(value)) {
			t.Errorf("Chord has incorrect value for key %v: got %v, want %v", key, string(storedValue), string(value))
		}
		if !exist {
			t.Errorf("Chord does not have the value for key %v", key)
		}
	}

	for i := 0; i < NumberNodes; i++ {
		node := nodes[ids[i]]
		node.StopWorking()
	}
	barrier.Wait()
}

func TestBasicReplication(t *testing.T) {
	NumberNodes := 10
	_, barrier, nodes, ids := StartUp("TestBasicReplication", NumberNodes)
	firstNodeId := ids[0]
	nodes[firstNodeId].Put(firstNodeId+1, []byte("test value"))
	time.Sleep((10 * WaitingTime) * time.Second)
	for i := 0; i < 2; i++ {
		nodes[ids[i]].StopWorking()
	}
	time.Sleep((10 * WaitingTime) * time.Second)
	value, exist := nodes[ids[2]].Get(firstNodeId + 1)
	if !exist {
		t.Errorf("Chord does not have the value for key %v", firstNodeId+1)
	}
	if exist && string(value) != "test value" {
		t.Errorf("Chord has incorrect value for key %v: got %v, want %v", firstNodeId+1, string(value), "test value")
	}
	for i := 2; i < NumberNodes; i++ {
		node := nodes[ids[i]]
		node.StopWorking()
	}
	barrier.Wait()
}

func TestReplication(t *testing.T) {
	NumberNodes := 10
	_, barrier, nodes, ids := StartUp("TestReplication", NumberNodes)
	firstNodeId := ids[0]
	nodes[firstNodeId].Put(firstNodeId+1, []byte("test value"))
	time.Sleep((10 * WaitingTime) * time.Second)
	for i := 0; i < 2; i++ {
		nodes[ids[i]].StopWorking()
	}
	time.Sleep((10 * WaitingTime) * time.Second)
	for i := 2; i < 4; i++ {
		nodes[ids[i]].StopWorking()
	}
	time.Sleep((10 * WaitingTime) * time.Second)
	value, exist := nodes[ids[4]].Get(firstNodeId + 1)
	if !exist {
		t.Errorf("Chord does not have the value for key %v", firstNodeId+1)
	}
	if exist && string(value) != "test value" {
		t.Errorf("Chord has incorrect value for key %v: got %v, want %v", firstNodeId+1, string(value), "test value")
	}
	for i := 4; i < NumberNodes; i++ {
		node := nodes[ids[i]]
		node.StopWorking()
	}
	barrier.Wait()
}
