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
	dict   map[Core.ChordHash]nodeDB
	logger common.Logger
}

func (db *DataBaseInMemory) GetNodesIds() []Core.ChordHash {
	var ids []Core.ChordHash
	for id := range db.dict {
		ids = append(ids, id)
	}
	return ids
}

func (db *DataBaseInMemory) GetNodeStateRPC(nodeId Core.ChordHash) string {
	nodeDB := db.dict[nodeId]
	return nodeDB.Node.GetState()
}

func NewDataBaseInMemory() *DataBaseInMemory {
	db := DataBaseInMemory{}
	db.logger = *common.NewLogger("DataBaseInMemory.txt")
	db.dict = make(map[Core.ChordHash]nodeDB)
	return &db
}

func (db *DataBaseInMemory) CreateRandomNode() {
	randomId := Core.GenerateRandomBinaryId()
	fmt.Println("Adding Node ", randomId)
	iString := strconv.Itoa(int(randomId))
	var server = NewServerInMemory(db, "Server"+iString)
	var client = NewClientInMemory(db, "Client"+iString)
	var monitor = MonitorHand.NewMonitorHand[ContactInMemory]("monitor" + iString)
	node := Core.NewBruteChord[ContactInMemory](server, client, monitor, randomId)
	db.addNode(node, server, client)
}

func (db *DataBaseInMemory) RemoveRandomNode() {
	nodeDB := Core.GetRandomFromDict(db.dict)
	db.RemoveNode(nodeDB.Node.GetId())
}

func (db *DataBaseInMemory) PutRandomDataRandomNode() {
	nodeDB := Core.GetRandomFromDict(db.dict)
	key := Core.GenerateRandomBinaryId()
	val := []byte{byte(key)}
	db.Put(nodeDB.Node.GetId(), key, val)
}

func (db *DataBaseInMemory) Put(nodeId Core.ChordHash, key Core.ChordHash, val []byte) {
	nodeDB := db.dict[nodeId]
	node := nodeDB.Node
	fmt.Printf("Going to put key %v with data %v using query node = %v \n", key, val, node.GetId())
	node.Put(key, val)
}

func (db *DataBaseInMemory) addNode(node *Core.BruteChord[ContactInMemory], server *ServerInMemory, client *ClientInMemory) {
	db.logger.WriteToFileOK("Adding node %v to database", node.GetId())
	db.dict[node.GetId()] = nodeDB{
		Node:   node,
		Server: server,
		Client: client,
	}
}

func (db *DataBaseInMemory) getServers() []*ServerInMemory {
	var servers []*ServerInMemory
	for _, node := range db.dict {
		servers = append(servers, node.Server)
	}
	return servers
}

func (db *DataBaseInMemory) getRandomNode() *Core.BruteChord[ContactInMemory] {
	nodeDB := Core.GetRandomFromDict(db.dict)
	return nodeDB.Node
}

func (db *DataBaseInMemory) GetNodes() []*Core.BruteChord[ContactInMemory] {
	var nodes []*Core.BruteChord[ContactInMemory]
	for _, node := range db.dict {
		nodes = append(nodes, node.Node)
	}
	return nodes
}

func (db *DataBaseInMemory) ChangeNodeState(nodeId Core.ChordHash, newState bool) {
	nodeDB := db.dict[nodeId]
	if newState {
		nodeDB.Node.SetWork(true)
	} else {
		nodeDB.Node.SetWork(false)
	}
}

func (db *DataBaseInMemory) RemoveNode(nodeId Core.ChordHash) {
	db.logger.WriteToFileOK("Removing node %v from database", nodeId)
	db.ChangeNodeState(nodeId, false)
	delete(db.dict, nodeId)
}
