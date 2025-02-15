package library

import (
	"iter"
	"maps"
)

type DataBaseInMemory struct {
	dict map[string]ServerInMemory
}

func NewDataBaseInMemory() *DataBaseInMemory {
	return &DataBaseInMemory{dict: make(map[string]ServerInMemory)}
}

func (db *DataBaseInMemory) AddServer(server *ServerInMemory) {
	db.dict[server.ServerId] = *server
}

func (db *DataBaseInMemory) GetServers() iter.Seq[ServerInMemory] {
	return maps.Values(db.dict)
}

type InMemoryContact struct {
	Id     string
	NodeId [NumberBits]uint8
}

func (i InMemoryContact) getNodeId() [NumberBits]uint8 {
	return i.NodeId
}

type ClientInMemory struct {
	DataBase *DataBaseInMemory
}

type ServerInMemory struct {
	DataBase             *DataBaseInMemory
	ServerId             string
	NodeId               [NumberBits]uint8
	ChannelCommunication chan Notification[InMemoryContact]
}

func (s *ServerInMemory) GetContact() InMemoryContact {
	return InMemoryContact{
		Id:     s.ServerId,
		NodeId: s.NodeId,
	}
}

func (s *ServerInMemory) SetData(channel chan Notification[InMemoryContact], Id [NumberBits]uint8) {
	s.ChannelCommunication = channel
	s.NodeId = Id
}

func (c *ClientInMemory) sendRequest(task ClientTask[InMemoryContact]) {
	database := c.DataBase
	for server := range database.GetServers() {
		for _, target := range task.Targets {
			if target.getNodeId() == server.GetContact().getNodeId() {
				server.ChannelCommunication <- task.Data
			}
		}
	}
}

func (c *ClientInMemory) sendRequestEveryone(data Notification[InMemoryContact]) {
	database := c.DataBase
	for server := range database.GetServers() {
		server.ChannelCommunication <- data
	}
}
