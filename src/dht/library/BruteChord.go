package library

import (
	"bittorrent/common"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type BruteChord[T Contact] struct {
	id                            [NumberBits]uint8    // Every Node should have an ID.
	successor                     T                    // A Contact, hidden because Set and Get Methods will be overloaded.
	lock                          sync.Mutex           // The Pointers shouldn't be updated concurrently.
	Monitor                       Monitor[T]           // To Keep Track of HeartBeats.
	NotificationChannelServerNode chan Notification[T] // A channel that will be intermediary between the Server and the Node.
	ServerChordCommunication      Server[T]            // A Server that will receive notifications from contacts of type T.
	ClientChordCommunication      Client[T]            // A Client that will send notifications to others nodes of type T.
	logger                        common.Logger        // To Log Everything The Node is doing.
	DeadContacts                  []T
}

func NewBruteChord[T Contact](serverChordCommunication Server[T], clientChordCommunication Client[T], monitor Monitor[T]) *BruteChord[T] {
	var node = BruteChord[T]{}
	node.id = GenerateRandomBinaryId()
	node.lock = sync.Mutex{}
	node.DeadContacts = make([]T, 0)
	node.logger = *common.NewLogger(strconv.Itoa(BinaryArrayToInt(node.id)) + ".txt")
	node.NotificationChannelServerNode = make(chan Notification[T])
	node.ServerChordCommunication = serverChordCommunication
	node.ServerChordCommunication.SetData(node.NotificationChannelServerNode, node.id)
	node.ClientChordCommunication = clientChordCommunication
	node.Monitor = monitor
	node.SetSuccessor(node.DefaultSuccessor())
	return &node
}

func (c *BruteChord[T]) GetId() [NumberBits]uint8 {
	return c.id
}

func (c *BruteChord[T]) GetSuccessor() T {
	c.lock.Lock()
	c.logger.WriteToFileOK(fmt.Sprintf("Calling GetSuccessor Method, returning %v", c.successor))
	successor := c.successor
	c.lock.Unlock()
	return successor
}

func (c *BruteChord[T]) SetSuccessor(candidate T) {
	curDate := time.Now()
	c.lock.Lock()
	c.Monitor.AddContact(candidate, curDate)
	c.logger.WriteToFileOK(fmt.Sprintf("Calling SetSuccessor Method, setting %v at date %v", candidate, curDate))
	c.successor = candidate
	c.lock.Unlock()
}

func (c *BruteChord[T]) DefaultSuccessor() T {
	c.logger.WriteToFileOK(fmt.Sprintf("Calling DefaultSuccessor Method, returning %v", c.ServerChordCommunication.GetContact()))
	return c.ServerChordCommunication.GetContact()
}

func (c *BruteChord[T]) GetContact() T {
	c.logger.WriteToFileOK(fmt.Sprintf("Calling GetContact Method, returning %v", c.ServerChordCommunication.GetContact()))
	return c.ServerChordCommunication.GetContact()
}

// BeginWorking Callers should use a Barrier because this is an infinite loop.
func (c *BruteChord[T]) BeginWorking() {
	c.logger.WriteToFileOK(fmt.Sprintf("Calling BeginWorking Method at data %v", time.Now()))
	c.cpu()
}

func (c *BruteChord[T]) cpu() {
	ticker := time.NewTicker(WaitingTime * time.Second)
	defer ticker.Stop() // Ensure the ticker stops when function exits
	for {
		select {
		case notification := <-c.NotificationChannelServerNode:
			c.logger.WriteToFileOK(fmt.Sprintf("Received Notification %v", reflect.TypeOf(notification)))
			// if this is called without go this stops working, but I don't know why, or how simulate the bug.
			if notification == nil {
				return
			}
			go notification.HandleNotification(c)
		case <-ticker.C:
			c.logger.WriteToFileOK(fmt.Sprintf("Cur time is %v, thus start doing Additional Things", time.Now()))
			go c.sendCheckPredecessor() // Stabilize Predecessor.
			go c.sendCheckAlive()       // Check if The Contacts I Have Are Alive.
			go c.killDead()             // Remove the Contacts that are Dead.
		}
	}
}

func (c *BruteChord[T]) sendCheckPredecessor() {
	c.logger.WriteToFileOK(fmt.Sprintf("Calling sendCheckPredecessor Method"))
	c.logger.WriteToFileOK(fmt.Sprintf("Sending AreYouMyPredecessor Notification to Everyone"))
	c.ClientChordCommunication.sendRequestEveryone(AreYouMyPredecessor[T]{
		Contact: c.GetContact(),
	})
}

func (c *BruteChord[T]) sendCheckAlive() {
	c.logger.WriteToFileOK(fmt.Sprintf("Calling sendCheckAlive Method"))
	c.logger.WriteToFileOK(fmt.Sprintf("Sending AreYouAliveNotification to %v", c.GetSuccessor()))
	c.ClientChordCommunication.sendRequest(ClientTask[T]{
		Targets: []T{c.GetSuccessor()},
		Data:    AreYouAliveNotification[T]{Contact: c.GetContact()},
	})
}

func (c *BruteChord[T]) killDead() {
	c.logger.WriteToFileOK(fmt.Sprintf("Calling killDead Method at time %v", time.Now()))
	successor := c.GetSuccessor()
	if !c.Monitor.CheckAlive(successor, 3*WaitingTime) {
		c.logger.WriteToFileOK(fmt.Sprintf("Successor %v is Dead", successor))
		c.DeadContacts = append(c.DeadContacts, successor)
		c.Monitor.DeleteContact(successor)
		c.SetSuccessor(c.DefaultSuccessor())
	} else {
		c.logger.WriteToFileOK(fmt.Sprintf("Successor %v is Alive", successor))
	}
}

func (c *BruteChord[T]) StopWorking() {
	c.NotificationChannelServerNode <- nil
}
