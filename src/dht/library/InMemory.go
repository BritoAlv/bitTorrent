package library

import (
	"bittorrent/common"
	"sync"
)

type DataBaseInMemory struct {
	lock   sync.Mutex
	dict   map[string]*ServerInMemory
	logger common.Logger
}

func NewDataBaseInMemory() *DataBaseInMemory {
	db := DataBaseInMemory{}
	db.logger = *common.NewLogger("DataBaseInMemory")
	db.lock = sync.Mutex{}
	db.dict = make(map[string]*ServerInMemory)
	return &db
}

func (db *DataBaseInMemory) AddServer(server *ServerInMemory) {
	db.logger.WriteToFileOK("Calling Adding server %s, waiting for lock", server.ServerId)
	db.lock.Lock()
	db.logger.WriteToFileOK("Lock is us, Adding server %s", server.ServerId)
	db.dict[server.ServerId] = server
	db.lock.Unlock()
}

func (db *DataBaseInMemory) GetServers() []*ServerInMemory {
	db.logger.WriteToFileOK("Calling GetServers, waiting for lock")
	db.lock.Lock()
	db.logger.WriteToFileOK("Lock is us, Creating a copy of the servers list")
	values := make([]*ServerInMemory, 0, len(db.dict))
	for _, value := range db.dict {
		values = append(values, value)
	}
	db.lock.Unlock()
	return values
}

type InMemoryContact struct {
	Id     string
	NodeId [NumberBits]uint8
}

func (i InMemoryContact) getNodeId() [NumberBits]uint8 {
	return i.NodeId
}
