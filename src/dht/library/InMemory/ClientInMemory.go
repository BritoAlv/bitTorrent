package InMemory

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord/Core"
	"reflect"
	"time"
)

type ClientInMemory struct {
	DataBase *DataBaseInMemory
	Logger   common.Logger
}

func NewClientInMemory(database *DataBaseInMemory, name string) *ClientInMemory {
	return &ClientInMemory{
		DataBase: database,
		Logger:   *common.NewLogger(name + ".txt"),
	}
}

func (c *ClientInMemory) SendRequest(task Core.ClientTask[ContactInMemory]) {
	c.Logger.WriteToFileOK("Calling SendRequest Method with the following task: %v at approximately date %v", task, time.Now())
	database := c.DataBase
	for _, server := range database.getServers() {
		for _, target := range task.Targets {
			if target.GetNodeId() == server.GetContact().GetNodeId() {
				c.Logger.WriteToFileOK("Sending data of type %v to specific server %v at approximately date %v", reflect.TypeOf(task.Data), server.GetContact(), time.Now())
				server.ChannelCommunication <- task.Data
			}
		}
	}
}

func (c *ClientInMemory) SendRequestEveryone(data Core.Notification[ContactInMemory]) {
	c.Logger.WriteToFileOK("Calling SendRequestEveryone Method with the following data: %v at approximately date %v", data, time.Now())
	database := c.DataBase
	for _, server := range database.getServers() {
		c.Logger.WriteToFileOK("Sending data to server %v", server.GetContact())
		server.ChannelCommunication <- data
	}
}
