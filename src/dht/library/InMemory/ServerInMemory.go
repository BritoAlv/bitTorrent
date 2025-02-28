package InMemory

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord/Core"
)

type ServerInMemory struct {
	DataBase             *DataBaseInMemory
	ServerId             string
	NodeId               Core.ChordHash
	ChannelCommunication chan Core.Notification[ContactInMemory]
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

func (s *ServerInMemory) GetContact() ContactInMemory {
	contact := ContactInMemory{
		Id:     s.ServerId,
		NodeId: s.NodeId,
	}
	s.Logger.WriteToFileOK("Calling GetContact Method returning %v", contact)
	return contact
}

func (s *ServerInMemory) SetData(channel chan Core.Notification[ContactInMemory], Id Core.ChordHash) {
	s.ChannelCommunication = channel
	s.NodeId = Id
}
