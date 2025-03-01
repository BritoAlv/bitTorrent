package InMemory

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord/Core"
	"math/rand"
	"testing"
	"time"
)

const NumberOfRuns = 2

func TestRunMultipleTimes(t *testing.T) {
	for i := 0; i < NumberOfRuns; i++ {
		t.Run("TestBasicChordBehaviourInitialization", TestBasicChordBehaviourInitialization)
		t.Run("TestBasicChordBehaviourStabilization", TestBasicChordBehaviourStabilization)
		t.Run("TestBasicPutGet", TestBasicPutGet)
		t.Run("TestBasicReplication", TestBasicReplication)
		t.Run("TestReplication", TestReplication)
		t.Run("TestTolerance", TestTolerance)
		t.Run("TestUpdateWithReplication", TestUpdateWithReplication)
	}
}

func StartUp(name string, NumberNodes int) (*DataBaseInMemory, map[Core.ChordHash]*Core.BruteChord[ContactInMemory], []Core.ChordHash) {
	common.SetLogDirectoryPath(name)
	var database = *NewDataBaseInMemory()
	for i := 0; i < NumberNodes; i++ {
		database.CreateRandomNode()
	}
	nodes := make(map[Core.ChordHash]*Core.BruteChord[ContactInMemory])
	ids := make([]Core.ChordHash, 0)
	for _, node := range database.GetNodes() {
		nodes[node.GetId()] = node
		ids = append(ids, node.GetId())
	}
	Core.Sort(ids)
	time.Sleep((10 * Core.WaitingTime) * time.Second)
	return &database, nodes, ids
}

/*
TestBasicChordBehaviourInitialization : N nodes are created simultaneously, eventually after stabilization occurs all the
nodes should have as its successor the next one in the ring.
*/
func TestBasicChordBehaviourInitialization(t *testing.T) {
	NumberNodes := 10
	_, nodes, ids := StartUp("TestBasicChordBehaviourInitialization", NumberNodes)
	for i := 0; i < NumberNodes; i++ {
		nodeId := ids[i]
		node := nodes[nodeId]
		successor := node.GetContact(1)
		successorId := successor.NodeId
		expectedSuccessorId := ids[(i+1)%NumberNodes]
		if successorId != expectedSuccessorId {
			t.Errorf("Node %v has successor %v, expected %v", nodeId, successorId, expectedSuccessorId)
		} else {
			t.Logf("Node %v has successor %v", nodeId, successorId)
		}
	}
}

// TestBasicChordBehaviourStabilization : N nodes are created simultaneously, some nodes randomly go down, and eventually go up, after stabilization occurs it should
// happen that all the nodes are successors are fine.
func TestBasicChordBehaviourStabilization(t *testing.T) {
	NumberNodes := 10
	database, nodes, ids := StartUp("TestBasicChordBehaviourStabilization", NumberNodes)
	down := make([]*Core.BruteChord[ContactInMemory], 0, NumberNodes)
	for i := 0; i < NumberNodes; i++ {
		if rand.Float32() <= 0.5 {
			nodeId := ids[i]
			node := nodes[nodeId]
			down = append(down, node)
			database.ChangeNodeState(nodeId, false)
		}
	}
	for _, node := range down {
		database.ChangeNodeState(node.GetId(), true)
	}
	time.Sleep((10 * Core.WaitingTime) * time.Second)
	for i := 0; i < NumberNodes; i++ {
		nodeId := ids[i]
		node := nodes[nodeId]
		successor := node.GetContact(1)
		successorId := successor.NodeId
		expectedSuccessorId := ids[(i+1)%NumberNodes]
		if successorId != expectedSuccessorId {
			t.Errorf("Node %v has successor %v, expected %v", nodeId, successorId, expectedSuccessorId)
		} else {
			t.Logf("Node %v has successor %v", nodeId, successorId)
		}
	}
}

func TestBasicPutGet(t *testing.T) {
	NumberNodes := 10
	database, nodes, _ := StartUp("TestBasicPutGet", NumberNodes)
	// Put data into a random node
	key := Core.GenerateRandomBinaryId()
	value := []byte("test value")
	randomNode := database.getRandomNode()
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
}

func TestBasicReplication(t *testing.T) {
	NumberNodes := 10
	db, nodes, ids := StartUp("TestBasicReplication", NumberNodes)
	firstNodeId := ids[0]
	nodes[firstNodeId].Put(firstNodeId+1, []byte("test value"))
	time.Sleep((10 * Core.WaitingTime) * time.Second)
	for i := 0; i < 2; i++ {
		db.ChangeNodeState(ids[i], false)
	}
	time.Sleep((10 * Core.WaitingTime) * time.Second)
	value, exist := nodes[ids[2]].Get(firstNodeId + 1)
	if !exist {
		t.Errorf("Chord does not have the value for key %v", firstNodeId+1)
	}
	if exist && string(value) != "test value" {
		t.Errorf("Chord has incorrect value for key %v: got %v, want %v", firstNodeId+1, string(value), "test value")
	}
}

func TestReplication(t *testing.T) {
	NumberNodes := 10
	db, nodes, ids := StartUp("TestReplication", NumberNodes)
	firstNodeId := ids[0]
	nodes[firstNodeId].Put(firstNodeId+1, []byte("test value"))
	time.Sleep((10 * Core.WaitingTime) * time.Second)
	for i := 0; i < 2; i++ {
		db.ChangeNodeState(ids[i], false)
	}
	time.Sleep((10 * Core.WaitingTime) * time.Second)
	for i := 2; i < 4; i++ {
		db.ChangeNodeState(ids[i], false)
	}
	time.Sleep((10 * Core.WaitingTime) * time.Second)
	value, exist := nodes[ids[4]].Get(firstNodeId + 1)
	if !exist {
		t.Errorf("Chord does not have the value for key %v", firstNodeId+1)
	}
	if exist && string(value) != "test value" {
		t.Errorf("Chord has incorrect value for key %v: got %v, want %v", firstNodeId+1, string(value), "test value")
	}
}

func TestTolerance(t *testing.T) {
	NumberNodes := 0
	Iteration := 10
	NumberDataLookUp := 10
	database, _, _ := StartUp("TestTolerance", NumberNodes)
	data := make(map[Core.ChordHash][]byte)
	for i := 0; i < NumberDataLookUp; i++ {
		randKey := Core.GenerateRandomBinaryId()
		data[randKey] = []byte{byte(randKey)}
	}
	for i := 0; i < Iteration; i++ {
		// choose two random nodes and stop them.
		t.Logf("Iteration %v", i)
		if len(database.GetNodes()) >= 3 {
			t.Logf("Removing nodes")
			for j := 0; j < 2; j++ {
				node1 := database.getRandomNode()
				database.RemoveNode(node1.GetId())
			}
			time.Sleep(4 * time.Second)
			node1 := database.getRandomNode()
			for key := range data {
				_, exist := node1.Get(key)
				if !exist {
					t.Fatalf("Chord does not have the value for key %v", key)
				} else {
					t.Logf("Query for %v found it", key)
				}
			}
		}
		if len(database.GetNodes()) <= 2 {
			t.Logf("Adding Nodes")
			for j := 0; j < 2; j++ {
				database.CreateRandomNode()
			}
		}
	}
}

func TestUpdateWithReplication(t *testing.T) {
	NumberNodes := 3
	database, _, _ := StartUp("TestUpdateWithReplication", NumberNodes)
	key := 10
	valueOld := []byte("old value")
	valueNew := []byte("new value")
	randomNode := database.getRandomNode()
	database.Put(randomNode.GetId(), Core.ChordHash(key), valueOld)
	time.Sleep(3 * time.Second)
	randomNode = database.getRandomNode()
	database.Put(randomNode.GetId(), Core.ChordHash(key), valueNew)
	time.Sleep(3 * time.Second)
	for _, node := range database.GetNodes() {
		value, exist := node.Get(Core.ChordHash(key))
		if !exist {
			t.Fatalf("Chord does not have the value for key %v", key)
		}
		if string(value) != string(valueNew) {
			t.Fatalf("Chord has incorrect value for key %v: got %v, want %v", key, string(value), string(valueNew))
		}
	}
	for i := 0; i < 2; i++ {
		database.RemoveRandomNode()
	}
	time.Sleep(8 * time.Second)
	for _, node := range database.GetNodes() {
		value, exist := node.Get(Core.ChordHash(key))
		if !exist {
			t.Fatalf("Chord does not have the value for key %v", key)
		}
		if string(value) != string(valueNew) {
			t.Fatalf("Chord has incorrect value for key %v: got %v, want %v", key, string(value), string(valueNew))
		}
	}
}
