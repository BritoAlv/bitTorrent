package library

import (
	"bittorrent/common"
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

func (c *ClientInMemory) sendRequest(task ClientTask[InMemoryContact]) {
	c.Logger.WriteToFileOK("Calling sendRequest Method with the following task: %v at approximately date %v", task, time.Now())
	database := c.DataBase
	for _, server := range database.GetServers() {
		for _, target := range task.Targets {
			if target.getNodeId() == server.GetContact().getNodeId() {
				c.Logger.WriteToFileOK("Sending data of type %v to specific server %v at approximately date %v", reflect.TypeOf(task.Data), server.GetContact(), time.Now())
				server.ChannelCommunication <- task.Data
			}
		}
	}
}

func (c *ClientInMemory) sendRequestEveryone(data Notification[InMemoryContact]) {
	c.Logger.WriteToFileOK("Calling sendRequestEveryone Method with the following data: %v at approximately date %v", data, time.Now())
	database := c.DataBase
	for _, server := range database.GetServers() {
		c.Logger.WriteToFileOK("Sending data to server %v", server.GetContact())
		server.ChannelCommunication <- data
	}
}
