package library

import (
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"
)

const NumberNodes = 20
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
	}
}

func StartUp(name string) (*DataBaseInMemory, *sync.WaitGroup) {
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
		go func() {
			barrier.Add(1)
			node.BeginWorking()
			barrier.Done()
		}()
	}
	return &database, &barrier
}

/*
TestBasicChordBehaviourInitialization : N nodes are created simultaneously, eventually after stabilization occurs all the
nodes should have as its successor the next one in the ring.
*/
func TestBasicChordBehaviourInitialization(t *testing.T) {
	database, barrier := StartUp("TestBasicChordBehaviourInitialization")
	nodes := make(map[ChordHash]*BruteChord[InMemoryContact])
	ids := make([]ChordHash, 0)
	for _, node := range database.GetNodes() {
		nodes[node.GetId()] = node
		ids = append(ids, node.GetId())
	}
	Sort(ids)
	time.Sleep((10 * WaitingTime) * time.Second)
	for i := 0; i < NumberNodes; i++ {
		nodeId := ids[i]
		node := nodes[nodeId]
		successor := node.GetSuccessor()
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
	database, barrier := StartUp("TestBasicChordBehaviourNoDead")
	nodes := database.GetNodes()
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
	database, barrier := StartUp("TestBasicChordBehaviourStabilization")
	nodes := make(map[ChordHash]*BruteChord[InMemoryContact])
	ids := make([]ChordHash, 0)
	for _, node := range database.GetNodes() {
		nodes[node.GetId()] = node
		ids = append(ids, node.GetId())
	}
	Sort(ids)
	time.Sleep((10 * WaitingTime) * time.Second)
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
		go func() {
			barrier.Add(1)
			node.BeginWorking()
			barrier.Done()
		}()
	}
	time.Sleep((10 * WaitingTime) * time.Second)
	for i := 0; i < NumberNodes; i++ {
		nodeId := ids[i]
		node := nodes[nodeId]
		successor := node.GetSuccessor()
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

// Do Some Put operations and check that the Data is inserted where it should be.
func TestBasicChordPut(t *testing.T) {

}
