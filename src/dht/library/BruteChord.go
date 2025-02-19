package library

import (
	"bittorrent/common"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type BruteChord[T Contact] struct {
	id                            [NumberBits]uint8    // Every Node should have an ID.
	successorList                 T                    // A Contact, hidden because Set and Get Methods will be overloaded.
	lock                          sync.Mutex           // The Pointers shouldn't be updated concurrently.
	Monitor                       Monitor[T]           // To Keep Track of HeartBeats.
	NotificationChannelServerNode chan Notification[T] // A channel that will be intermediary between the Server and the Node.
	ServerChordCommunication      Server[T]            // A Server that will receive notifications from contacts of type T.
	ClientChordCommunication      Client[T]            // A Client that will send notifications to others nodes of type T.
	logger                        common.Logger        // To Log Everything The Node is doing.
	DeadContacts                  []T                  // This is only for testing purposes.
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
	c.logger.WriteToFileOK("Calling GetSuccessor Method, returning %v", c.successorList)
	successor := c.successorList
	c.lock.Unlock()
	return successor
}

func (c *BruteChord[T]) SetSuccessor(candidate T) {
	curDate := time.Now()
	c.lock.Lock()
	c.Monitor.AddContact(candidate, curDate)
	c.logger.WriteToFileOK("Calling SetSuccessor Method, setting %v at date %v", candidate, curDate)
	c.successorList = candidate
	c.lock.Unlock()
}

func (c *BruteChord[T]) DefaultSuccessor() T {
	c.logger.WriteToFileOK("Calling DefaultSuccessor Method, returning %v", c.ServerChordCommunication.GetContact())
	return c.ServerChordCommunication.GetContact()
}

func (c *BruteChord[T]) GetContact() T {
	c.logger.WriteToFileOK("Calling GetContact Method, returning %v", c.ServerChordCommunication.GetContact())
	return c.ServerChordCommunication.GetContact()
}

// BeginWorking Callers should use a Barrier because this is an infinite loop.
func (c *BruteChord[T]) BeginWorking() {
	c.logger.WriteToFileOK("Calling BeginWorking Method at data %v", time.Now())
	c.cpu()
}

func (c *BruteChord[T]) cpu() {
	ticker := time.NewTicker(WaitingTime * time.Second)
	defer ticker.Stop() // Ensure the ticker stops when function exits
	for {
		select {
		case notification := <-c.NotificationChannelServerNode:
			c.logger.WriteToFileOK("Received Notification %v", reflect.TypeOf(notification))
			// if this is called without go this stops working, but I don't know why, or how simulate the bug.
			if notification == nil {
				return
			}
			go notification.HandleNotification(c)
		case <-ticker.C:
			c.logger.WriteToFileOK("Cur time is %v, thus start doing Additional Things", time.Now())
			go c.sendCheckPredecessor() // Stabilize Predecessor.
			go c.sendCheckAlive()       // Check if The Contacts I Have Are Alive.
			go c.killDead()             // Remove the Contacts that are Dead.
		}
	}
}

func (c *BruteChord[T]) sendCheckPredecessor() {
	c.logger.WriteToFileOK("Calling sendCheckPredecessor Method")
	c.logger.WriteToFileOK("Sending AreYouMyPredecessor Notification to Everyone")
	c.ClientChordCommunication.sendRequestEveryone(AreYouMyPredecessor[T]{
		Contact: c.GetContact(),
	})
}

func (c *BruteChord[T]) sendCheckAlive() {
	c.logger.WriteToFileOK("Calling sendCheckAlive Method")
	c.logger.WriteToFileOK("Sending AreYouAliveNotification to %v", c.GetSuccessor())
	c.ClientChordCommunication.sendRequest(ClientTask[T]{
		Targets: []T{c.GetSuccessor()},
		Data:    AreYouAliveNotification[T]{Contact: c.GetContact()},
	})
}

func (c *BruteChord[T]) killDead() {
	c.logger.WriteToFileOK("Calling killDead Method at time %v", time.Now())
	successor := c.GetSuccessor()
	if !c.Monitor.CheckAlive(successor, 3*WaitingTime) {
		c.logger.WriteToFileOK("Successor %v is Dead", successor)
		c.DeadContacts = append(c.DeadContacts, successor)
		c.Monitor.DeleteContact(successor)
		c.SetSuccessor(c.DefaultSuccessor())
	} else {
		c.logger.WriteToFileOK("Successor %v is Alive", successor)
	}
}

func (c *BruteChord[T]) StopWorking() {
	c.NotificationChannelServerNode <- nil
}
