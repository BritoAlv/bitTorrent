package library

import (
	"bittorrent/common"
	"sync"
)

type NodeDB struct {
	Node   *BruteChord[InMemoryContact]
	Server *ServerInMemory
	Client *ClientInMemory
}

type DataBaseInMemory struct {
	lock   sync.Mutex
	dict   map[ChordHash]NodeDB
	logger common.Logger
}

func NewDataBaseInMemory() *DataBaseInMemory {
	db := DataBaseInMemory{}
	db.logger = *common.NewLogger("DataBaseInMemory")
	db.lock = sync.Mutex{}
	db.dict = make(map[ChordHash]NodeDB)
	return &db
}

func (db *DataBaseInMemory) AddNode(node *BruteChord[InMemoryContact], server *ServerInMemory, client *ClientInMemory) {
	db.lock.Lock()
	db.logger.WriteToFileOK("Adding node %v to database", node.GetId())
	db.dict[node.id] = NodeDB{
		Node:   node,
		Server: server,
		Client: client,
	}
	db.lock.Unlock()
}

func (db *DataBaseInMemory) GetServers() []*ServerInMemory {
	db.lock.Lock()
	defer db.lock.Unlock()
	var servers []*ServerInMemory
	for _, node := range db.dict {
		servers = append(servers, node.Server)
	}
	return servers
}

func (db *DataBaseInMemory) GetNodes() []*BruteChord[InMemoryContact] {
	db.lock.Lock()
	defer db.lock.Unlock()
	var nodes []*BruteChord[InMemoryContact]
	for _, node := range db.dict {
		nodes = append(nodes, node.Node)
	}
	return nodes
}

func (db *DataBaseInMemory) RemoveNode(node *BruteChord[InMemoryContact]) {
	db.lock.Lock()
	db.logger.WriteToFileOK("Removing node %v from database", node.GetId())
	node.StopWorking()
	delete(db.dict, node.id)
	db.lock.Unlock()
}

type InMemoryContact struct {
	Id     string
	NodeId ChordHash
}

func (i InMemoryContact) getNodeId() ChordHash {
	return i.NodeId
}
