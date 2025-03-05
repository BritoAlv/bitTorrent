package InMemory

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord/Core"
	"bittorrent/dht/library/MonitorHand"
	"fmt"
	"strconv"
)

type nodeDB struct {
	Node   *Core.BruteChord[ContactInMemory]
	Server *ServerInMemory
	Client *ClientInMemory
}

type DataBaseInMemory struct { // Not Meant To Be Async.
	dict   Core.SafeMap[Core.ChordHash, nodeDB]
	logger common.Logger
}

func (db *DataBaseInMemory) GetActiveNodesIds() []Core.ChordHash {
	var ids []Core.ChordHash
	for _, id := range db.dict.GetKeys() {
		ids = append(ids, id)
	}
	return ids
}

func (db *DataBaseInMemory) GetNodeStateRPC(nodeId Core.ChordHash) string {
	nodeDB, _ := db.dict.Get(nodeId)
	nodeState := nodeDB.Node.GetState()
	return nodeState.String()
}

func NewDataBaseInMemory() *DataBaseInMemory {
	db := DataBaseInMemory{}
	db.logger = *common.NewLogger("DataBaseInMemory.txt")
	db.dict = *Core.NewSafeMap[Core.ChordHash, nodeDB](make(map[Core.ChordHash]nodeDB))
	return &db
}

func (db *DataBaseInMemory) CreateRandomNode() Core.ChordHash {
	randomId := Core.GenerateRandomBinaryId()
	fmt.Println("Adding Node ", randomId)
	iString := strconv.Itoa(int(randomId))
	var server = NewServerInMemory(db, "Server"+iString)
	var client = NewClientInMemory(db, "Client"+iString)
	var monitor = MonitorHand.NewMonitorHand[ContactInMemory]("monitor" + iString)
	node := Core.NewBruteChord[ContactInMemory](server, client, monitor, randomId)
	db.addNode(node, server, client)
	return randomId
}

func (db *DataBaseInMemory) RemoveRandomNode(n int) []Core.ChordHash {
	nodes := db.getRandomNode(n)
	ids := make([]Core.ChordHash, n)
	for i, node := range nodes {
		id := node.GetId()
		db.RemoveNode(id)
		ids[i] = id
	}
	return ids
}

func (db *DataBaseInMemory) PutRandomDataRandomNode() {
	node := db.getRandomNode(1)
	key := Core.GenerateRandomBinaryId()
	val := []byte{byte(key)}
	node[0].Put(key, val)
}

func (db *DataBaseInMemory) Put(nodeId Core.ChordHash, key Core.ChordHash, val []byte) {
	nodeDict, _ := db.dict.Get(nodeId)
	node := nodeDict.Node
	fmt.Printf("Going to put key %v with data %v using query node = %v \n", key, val, node.GetId())
	node.Put(key, val)
}

func (db *DataBaseInMemory) addNode(node *Core.BruteChord[ContactInMemory], server *ServerInMemory, client *ClientInMemory) {
	db.logger.WriteToFileOK("Adding node %v to database", node.GetId())
	db.dict.Set(node.GetId(), nodeDB{
		Node:   node,
		Server: server,
		Client: client,
	})
}

func (db *DataBaseInMemory) getServers() []*ServerInMemory {
	var servers []*ServerInMemory
	for _, node := range db.dict.GetValues() {
		servers = append(servers, node.Server)
	}
	return servers
}

func (db *DataBaseInMemory) getRandomNode(n int) []*Core.BruteChord[ContactInMemory] {
	nodes := db.dict.GetValues()
	result := make([]*Core.BruteChord[ContactInMemory], n)
	for i := 0; i < n; i++ {
		result[i] = nodes[i].Node
	}
	return result
}

func (db *DataBaseInMemory) GetNodes() []*Core.BruteChord[ContactInMemory] {
	var nodes []*Core.BruteChord[ContactInMemory]
	for _, node := range db.dict.GetValues() {
		nodes = append(nodes, node.Node)
	}
	return nodes
}

func (db *DataBaseInMemory) ChangeNodeState(nodeId Core.ChordHash, newState bool) {
	node, _ := db.dict.Get(nodeId)
	if newState {
		node.Node.SetWork(true)
	} else {
		node.Node.SetWork(false)
	}
}

func (db *DataBaseInMemory) RemoveNode(nodeId Core.ChordHash) {
	db.logger.WriteToFileOK("Removing node %v from database", nodeId)
	db.ChangeNodeState(nodeId, false)
	db.dict.Delete(nodeId)
}
