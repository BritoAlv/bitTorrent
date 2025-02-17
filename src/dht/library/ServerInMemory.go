package library

import (
	"bittorrent/common"
	"fmt"
)

type ServerInMemory struct {
	DataBase             *DataBaseInMemory
	ServerId             string
	NodeId               [NumberBits]uint8
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
	s.Logger.WriteToFileOK(fmt.Sprintf("Calling GetContact Method returning %v", contact))
	return contact
}

func (s *ServerInMemory) SetData(channel chan Notification[InMemoryContact], Id [NumberBits]uint8) {
	s.ChannelCommunication = channel
	s.NodeId = Id
}
