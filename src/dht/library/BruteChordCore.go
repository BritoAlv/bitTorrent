package library

import (
	"bittorrent/common"
	"reflect"
	"sync"
	"time"
)

type ContactWithData[T Contact] struct {
	Contact T
	Data    Store
}

type Confirmations struct {
	Confirmation bool
	Value        []byte // this have to be interpreted by the caller by now.
}

type BruteChord[T Contact] struct {
	id                            ChordHash // The ID of the Node.
	Info                          [3]ContactWithData[T]
	predecessorRef                T                       // Keep Track of the Predecessor, due to replication.
	lock                          sync.Mutex              // The Pointers shouldn't be updated concurrently.
	Monitor                       Monitor[T]              // To Keep Track of HeartBeats.
	NotificationChannelServerNode chan Notification[T]    // A channel that will be intermediary between the Server and the Node.
	ServerChordCommunication      Server[T]               // A Server that will receive notifications from contacts of type T.
	ClientChordCommunication      Client[T]               // A Client that will send notifications to others nodes of type T.
	logger                        common.Logger           // To Log Everything The Node is doing.
	DeadContacts                  []T                     // This is only for testing purposes.
	pendingResponses              map[int64]Confirmations // I need to check if someone answered the request I sent.
}

func NewBruteChord[T Contact](serverChordCommunication Server[T], clientChordCommunication Client[T], monitor Monitor[T], id ChordHash) *BruteChord[T] {
	var node = BruteChord[T]{}
	node.id = id
	node.lock = sync.Mutex{}
	node.DeadContacts = make([]T, 0)
	node.logger = *common.NewLogger(ToString(node.id) + ".txt")
	node.NotificationChannelServerNode = make(chan Notification[T])
	node.ServerChordCommunication = serverChordCommunication
	node.ServerChordCommunication.SetData(node.NotificationChannelServerNode, node.id)
	node.ClientChordCommunication = clientChordCommunication
	node.Monitor = monitor
	node.pendingResponses = make(map[int64]Confirmations)
	node.SetContact(node.DefaultSuccessor(), -1)
	node.Info[0].Contact = serverChordCommunication.GetContact()
	for i := 0; i < 3; i++ {
		node.Info[i].Data = make(Store)
	}
	node.SetContact(node.DefaultSuccessor(), 1)
	node.SetContact(node.DefaultSuccessor(), 2)
	return &node
}

func (c *BruteChord[T]) GetId() ChordHash {
	return c.id
}

func (c *BruteChord[T]) copyValue(value []byte) []byte {
	data := make([]byte, len(value), len(value))
	copy(data, value)
	return data
}

func (c *BruteChord[T]) copyStore(store Store) Store {
	data := make(map[ChordHash][]byte)
	for key, value := range store {
		data[key] = make([]byte, len(value), len(value))
		copy(data[key], value)
	}
	return data
}

func (c *BruteChord[T]) GetData(key ChordHash, index int) []byte {
	c.lock.Lock()
	defer c.lock.Unlock()
	data := c.copyValue(c.Info[index].Data[key])
	return data
}

func (c *BruteChord[T]) GetAllData(index int) Store {
	c.lock.Lock()
	defer c.lock.Unlock()
	data := c.copyStore(c.Info[index].Data)
	return data
}

func (c *BruteChord[T]) SetData(key ChordHash, value []byte, index int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Info[index].Data[key] = value
}

func (c *BruteChord[T]) SetPendingResponse(taskId int64, confirmation Confirmations) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.pendingResponses[taskId] = confirmation
}

func (c *BruteChord[T]) AddNewData(data Store, index int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Info[index].Data = data
}

func (c *BruteChord[T]) GetPendingResponse(taskId int64) Confirmations {
	c.lock.Lock()
	defer c.lock.Unlock()
	confirmation, _ := c.pendingResponses[taskId]
	return confirmation
}

func (c *BruteChord[T]) SetContact(candidate T, index int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if index == -1 {
		c.predecessorRef = candidate
	} else {
		if c.Info[index].Contact.getNodeId() != candidate.getNodeId() {
			c.logger.WriteToFileOK("Because I'll have a new successor, I'll have to release the replicas I was holding for that old successor")
			for key, data := range c.Info[index].Data {
				go c.Put(key, data)
			}
			c.Info[index] = ContactWithData[T]{Contact: candidate, Data: make(Store)}
		}
	}
	curDate := time.Now()
	c.Monitor.AddContact(candidate, curDate)
}

func (c *BruteChord[T]) DefaultSuccessor() T {
	contact := c.Info[0].Contact
	c.logger.WriteToFileOK("Calling DefaultSuccessor Method, returning %v", contact)
	return contact
}

func (c *BruteChord[T]) GetContact(index int) T {
	c.lock.Lock()
	defer c.lock.Unlock()
	var contact T
	if index == -1 {
		contact = c.predecessorRef
	} else {
		contact = c.Info[index].Contact
	}
	return contact
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
