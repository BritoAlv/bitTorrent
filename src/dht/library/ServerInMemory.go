package library

import (
	"bittorrent/common"
)

type ServerInMemory struct {
	DataBase             *DataBaseInMemory
	ServerId             string
	NodeId               ChordHash
	ChannelCommunication chan Notification[InMemoryContact]
	Logger               common.Logger
}

func NewServerInMemory(database *DataBaseInMemory, name string) *ServerInMemory {
	return &ServerInMemory{
		DataBase:             database,
		ServerId:             name,
		ChannelCommunication: nil,
		Logger:               *common.NewLogger(name + ".txt"),
	}
}

func (s *ServerInMemory) GetContact() InMemoryContact {
	contact := InMemoryContact{
		Id:     s.ServerId,
		NodeId: s.NodeId,
	}
	s.Logger.WriteToFileOK("Calling GetContact Method returning %v", contact)
	return contact
}

func (s *ServerInMemory) SetData(channel chan Notification[InMemoryContact], Id ChordHash) {
	s.ChannelCommunication = channel
	s.NodeId = Id
}
