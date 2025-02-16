package library

import (
	"iter"
	"maps"
)

type DataBaseInMemory struct {
	dict map[string]*ServerInMemory
}

func NewDataBaseInMemory() *DataBaseInMemory {
	return &DataBaseInMemory{dict: make(map[string]*ServerInMemory)}
}

func (db *DataBaseInMemory) AddServer(server *ServerInMemory) {
	db.dict[server.ServerId] = server
}

func (db *DataBaseInMemory) GetServers() iter.Seq[*ServerInMemory] {
	return maps.Values(db.dict)
}

type InMemoryContact struct {
	Id     string
	NodeId [NumberBits]uint8
}

func (i InMemoryContact) getNodeId() [NumberBits]uint8 {
	return i.NodeId
}
