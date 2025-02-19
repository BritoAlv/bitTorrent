package library

import (
	"bittorrent/common"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"
)

const NumberNodes = 30
const NumberOfRuns = 2

func SetLogDirectoryPath(name string) {
	common.LogsPath = "./logs/" + name + strconv.Itoa(time.Now().Nanosecond()) + "/"
}

func TestRunMultipleTimes(t *testing.T) {
	for i := 0; i < NumberOfRuns; i++ {
		t.Run("TestBasicChordBehaviourInitialization", TestBasicChordBehaviourInitialization)
		t.Run("TestBasicChordBehaviourNoDead", TestBasicChordBehaviourNoDead)
		t.Run("TestBasicChordBehaviourStabilization", TestBasicChordBehaviourStabilization)
	}
}

/*
TestBasicChordBehaviourInitialization : N nodes are created simultaneously, eventually after stabilization occurs all the
nodes should have as its successorList the next one in the ring.
*/
func TestBasicChordBehaviourInitialization(t *testing.T) {
	SetLogDirectoryPath("TestBasicChordBehaviourInitialization")
	var database = *NewDataBaseInMemory()
	var ids = make([]int, 0, NumberNodes)
	var nodes = make(map[int]*BruteChord[InMemoryContact])
	var barrier = sync.WaitGroup{}
	for i := 0; i < NumberNodes; i++ {
		iString := strconv.Itoa(i)
		var server = NewServerInMemory(&database, "Server"+iString)
		var client = NewClientInMemory(&database, "Client"+iString)
		database.AddServer(server)
		node := NewBruteChord[InMemoryContact](server, client, NewMonitorHand[InMemoryContact]("Monitor"+iString))
		intNodeId := BinaryArrayToInt(node.GetId())
		nodes[intNodeId] = node
		ids = append(ids, BinaryArrayToInt(node.GetId()))
		go func() {
			barrier.Add(1)
			node.BeginWorking()
			barrier.Done()
		}()
	}
	sort.Ints(ids)
	time.Sleep((3 * WaitingTime) * time.Second)
	for i := 0; i < NumberNodes; i++ {
		nodeId := ids[i]
		node := nodes[nodeId]
		successor := node.GetSuccessor()
		successorId := BinaryArrayToInt(successor.NodeId)
		expectedSuccessorId := ids[(i+1)%NumberNodes]
		if successorId != expectedSuccessorId {
			t.Errorf("Node %v has successorList %v, expected %v", nodeId, successorId, expectedSuccessorId)
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
	SetLogDirectoryPath("TestBasicChordBehaviourNoDead")
	var database = *NewDataBaseInMemory()
	var nodes = make(map[int]*BruteChord[InMemoryContact])
	var barrier = sync.WaitGroup{}
	for i := 0; i < NumberNodes; i++ {
		iString := strconv.Itoa(i)
		var server = NewServerInMemory(&database, "Server"+iString)
		var client = NewClientInMemory(&database, "Client"+iString)
		database.AddServer(server)
		node := NewBruteChord[InMemoryContact](server, client, NewMonitorHand[InMemoryContact]("Monitor"+iString))
		intNodeId := BinaryArrayToInt(node.GetId())
		nodes[intNodeId] = node
		go func() {
			barrier.Add(1)
			node.BeginWorking()
			barrier.Done()
		}()
	}
	time.Sleep((3 * WaitingTime) * time.Second)
	for _, node := range nodes {
		if len(node.DeadContacts) > 0 {
			t.Errorf("Node %v has the dead contacts ", BinaryArrayToInt(node.GetId()))
		}
	}
	for _, node := range nodes {
		node.StopWorking()
	}
	barrier.Wait()
}

// TestBasicChordBehaviourStabilization : N nodes are created simultaneously, some nodes randomly go down, and eventually go up, after stabilization occurs it should
// happen that all the nodes are alive.
func TestBasicChordBehaviourStabilization(t *testing.T) {
	SetLogDirectoryPath("TestBasicChordBehaviourStabilization")
	var database = *NewDataBaseInMemory()
	var ids = make([]int, 0, NumberNodes)
	var nodes = make(map[int]*BruteChord[InMemoryContact])
	var barrier = sync.WaitGroup{}
	for i := 0; i < NumberNodes; i++ {
		iString := strconv.Itoa(i)
		var server = NewServerInMemory(&database, "Server"+iString)
		var client = NewClientInMemory(&database, "Client"+iString)
		database.AddServer(server)
		node := NewBruteChord[InMemoryContact](server, client, NewMonitorHand[InMemoryContact]("Monitor"+iString))
		intNodeId := BinaryArrayToInt(node.GetId())
		nodes[intNodeId] = node
		ids = append(ids, BinaryArrayToInt(node.GetId()))
		go func() {
			barrier.Add(1)
			node.BeginWorking()
			barrier.Done()
		}()
	}
	sort.Ints(ids)
	time.Sleep((3 * WaitingTime) * time.Second)
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
	time.Sleep((3 * WaitingTime) * time.Second)
	for i := 0; i < NumberNodes; i++ {
		nodeId := ids[i]
		node := nodes[nodeId]
		successor := node.GetSuccessor()
		successorId := BinaryArrayToInt(successor.NodeId)
		expectedSuccessorId := ids[(i+1)%NumberNodes]
		if successorId != expectedSuccessorId {
			t.Errorf("Node %v has successorList %v, expected %v", nodeId, successorId, expectedSuccessorId)
		}
	}
	for i := 0; i < NumberNodes; i++ {
		node := nodes[ids[i]]
		node.StopWorking()
	}
	barrier.Wait()
}
