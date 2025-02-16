package library

import (
	"bittorrent/common"
	"fmt"
	"reflect"
	"time"
)

type BruteChord[T Contact] struct {
	Id                            [NumberBits]uint8    // Every Node should have an ID.
	successor                     T                    // A Contact, hidden because Set and Get Methods will be overloaded.
	Monitor                       Monitor[T]           // To Keep Track of HeartBeats.
	NotificationChannelServerNode chan Notification[T] // A channel that will be intermediary between the Server and the Node.
	ServerChordCommunication      Server[T]            // A Server that will receive notifications from contacts of type T.
	ClientChordCommunication      Client[T]            // A Client that will send notifications to others nodes of type T.
	Logger                        common.Logger        // To Log Everything The Node is doing.
}

func NewBruteChord[T Contact](serverChordCommunication Server[T], clientChordCommunication Client[T], monitor Monitor[T]) *BruteChord[T] {
	var node = BruteChord[T]{}
	node.Id = GenerateRandomBinaryId()
	node.Logger = *common.NewLogger(ConvertStr(node.Id) + ".txt")
	node.NotificationChannelServerNode = make(chan Notification[T])
	node.ServerChordCommunication = serverChordCommunication
	node.ServerChordCommunication.SetData(node.NotificationChannelServerNode, node.Id)
	node.ClientChordCommunication = clientChordCommunication
	node.Monitor = monitor
	node.SetSuccessor(node.DefaultSuccessor())
	return &node
}

func (c *BruteChord[T]) GetSuccessor() T {
	c.Logger.WriteToFileOK(fmt.Sprintf("Calling GetSuccessor Method, returning %v", c.successor))
	return c.successor
}

func (c *BruteChord[T]) SetSuccessor(candidate T) {
	curDate := time.Now()
	c.Monitor.AddContact(candidate, curDate)
	c.Logger.WriteToFileOK(fmt.Sprintf("Calling SetSuccessor Method, setting %v at date %v", candidate, curDate))
	c.successor = candidate
}

func (c *BruteChord[T]) DefaultSuccessor() T {
	c.Logger.WriteToFileOK(fmt.Sprintf("Calling DefaultSuccessor Method, returning %v", c.ServerChordCommunication.GetContact()))
	return c.ServerChordCommunication.GetContact()
}

func (c *BruteChord[T]) GetContact() T {
	c.Logger.WriteToFileOK(fmt.Sprintf("Calling GetContact Method, returning %v", c.ServerChordCommunication.GetContact()))
	return c.ServerChordCommunication.GetContact()
}

// BeginWorking Callers should use a Barrier because this is an infinite loop.
func (c *BruteChord[T]) BeginWorking() {
	c.Logger.WriteToFileOK(fmt.Sprintf("Calling BeginWorking Method at data %v", time.Now()))
	c.cpu()
}

func (c *BruteChord[T]) cpu() {
	ticker := time.NewTicker(WaitingTime * time.Second)
	defer ticker.Stop() // Ensure the ticker stops when function exits
	for {
		select {
		case notification := <-c.NotificationChannelServerNode:
			c.Logger.WriteToFileOK(fmt.Sprintf("Received Notification %v", reflect.TypeOf(notification)))
			notification.HandleNotification(c)
		case <-ticker.C:
			c.Logger.WriteToFileOK(fmt.Sprintf("Cur time is %v, thus start doing Additional Things", time.Now()))
			go c.sendCheckPredecessor() // Stabilize Predecessor.
			go c.sendCheckAlive()       // Check if The Contacts I Have Are Alive.
			go c.killDead()             // Remove the Contacts that are Dead.
		}
	}
}

func (c *BruteChord[T]) sendCheckPredecessor() {
	c.Logger.WriteToFileOK(fmt.Sprintf("Calling sendCheckPredecessor Method"))
	c.Logger.WriteToFileOK(fmt.Sprintf("Sending AreYouMyPredecessor Notification to Everyone"))
	c.ClientChordCommunication.sendRequestEveryone(AreYouMyPredecessor[T]{
		Contact: c.GetContact(),
	})
}

func (c *BruteChord[T]) sendCheckAlive() {
	c.Logger.WriteToFileOK(fmt.Sprintf("Calling sendCheckAlive Method"))
	c.Logger.WriteToFileOK(fmt.Sprintf("Sending AreYouAliveNotification to %v", c.GetSuccessor()))
	c.ClientChordCommunication.sendRequest(ClientTask[T]{
		Targets: []T{c.GetSuccessor()},
		Data:    AreYouAliveNotification[T]{Contact: c.GetContact()},
	})
}

func (c *BruteChord[T]) killDead() {
	c.Logger.WriteToFileOK(fmt.Sprintf("Calling killDead Method at time %v", time.Now()))
	successor := c.GetSuccessor()
	if !c.Monitor.CheckAlive(successor, 2*WaitingTime) {
		c.Logger.WriteToFileOK(fmt.Sprintf("Successor %v is Dead", successor))
		c.Monitor.DeleteContact(successor)
		c.SetSuccessor(c.DefaultSuccessor())
	} else {
		c.Logger.WriteToFileOK(fmt.Sprintf("Successor %v is Alive", successor))
	}
}
