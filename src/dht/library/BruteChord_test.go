package library

import (
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"
)

const NumberNodes = 10

/*
TestBasicChordBehaviourStabilization : N nodes are created simultaneously, eventually after stabilization occurs all the
nodes should have as its successor the next one in the circle.
*/
func TestBasicChordBehaviourStabilization(t *testing.T) {
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
	time.Sleep((NumberNodes * WaitingTime) * time.Second)
	for i := 0; i < NumberNodes; i++ {
		nodeId := ids[i]
		node := nodes[nodeId]
		successor := node.GetSuccessor()
		successorId := BinaryArrayToInt(successor.NodeId)
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
	time.Sleep((NumberNodes * WaitingTime) * time.Second)
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
