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
	id                            ChordHash               // Every Node should have an ID.
	successor                     ContactWithData[T]      // Contact of its successors, hidden because set and get methods will be overloaded.
	successorSuccessor            ContactWithData[T]      // I need to keep track of this to replicate data in case of network changes.
	predecessorRef                T                       // Keep Track of the Predecessor, due to replication.
	lock                          sync.Mutex              // The Pointers shouldn't be updated concurrently.
	Monitor                       Monitor[T]              // To Keep Track of HeartBeats.
	NotificationChannelServerNode chan Notification[T]    // A channel that will be intermediary between the Server and the Node.
	ServerChordCommunication      Server[T]               // A Server that will receive notifications from contacts of type T.
	ClientChordCommunication      Client[T]               // A Client that will send notifications to others nodes of type T.
	logger                        common.Logger           // To Log Everything The Node is doing.
	DeadContacts                  []T                     // This is only for testing purposes.
	myData                        Store                   // My Data That I'm actually Responsible for.
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
	node.successor.Data = make(Store)
	node.successorSuccessor.Data = make(Store)
	node.SetSuccessor(node.DefaultSuccessor())
	node.SetSuccessorSuccessor(node.DefaultSuccessor())
	node.SetPredecessor(node.DefaultSuccessor())
	node.pendingResponses = make(map[int64]Confirmations)
	node.myData = make(Store)
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

func (c *BruteChord[T]) GetData(key ChordHash) []byte {
	c.lock.Lock()
	defer c.lock.Unlock()
	data := c.copyValue(c.myData[key])
	return data
}

func (c *BruteChord[T]) GetAllOwnData() Store {
	c.lock.Lock()
	defer c.lock.Unlock()
	data := c.copyStore(c.myData)
	return data
}

func (c *BruteChord[T]) GetSuccessorReplicatedData() Store {
	c.lock.Lock()
	defer c.lock.Unlock()
	data := c.copyStore(c.successor.Data)
	return data
}

func (c *BruteChord[T]) GetSuccessorSuccessorReplicatedData() Store {
	c.lock.Lock()
	defer c.lock.Unlock()
	data := c.copyStore(c.successorSuccessor.Data)
	return data
}

func (c *BruteChord[T]) SetData(key ChordHash, value []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.myData[key] = value
}

func (c *BruteChord[T]) GetSuccessor() T {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.logger.WriteToFileOK("Calling GetSuccessor Method, returning %v", c.successor)
	successor := c.successor.Contact
	return successor
}

func (c *BruteChord[T]) SetPendingResponse(taskId int64, confirmation Confirmations) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.pendingResponses[taskId] = confirmation
}

func (c *BruteChord[T]) ReplaceSuccessorData(data Store) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.successor.Data = c.copyStore(data)
}

func (c *BruteChord[T]) ReplaceSuccessorSuccessorData(data Store) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.successorSuccessor.Data = c.copyStore(data)
}

func (c *BruteChord[T]) GetPendingResponse(taskId int64) Confirmations {
	c.lock.Lock()
	defer c.lock.Unlock()
	confirmation, _ := c.pendingResponses[taskId]
	return confirmation
}

func (c *BruteChord[T]) releaseSuccessorReplica() {
	c.lock.Lock()
	defer c.lock.Unlock()
	for key, data := range c.successor.Data {
		c.myData[key] = data
	}
	c.successor.Data = make(Store)
}

func (c *BruteChord[T]) releaseSuccessorSuccessorReplica() {
	c.lock.Lock()
	defer c.lock.Unlock()
	for key, data := range c.successorSuccessor.Data {
		c.myData[key] = data
	}
	c.successor.Data = make(Store)
}

func (c *BruteChord[T]) SetSuccessor(candidate T) {
	if c.successor.Contact.getNodeId() != candidate.getNodeId() {
		c.logger.WriteToFileOK("Because I'll have a new successor, I'll have to release the replicas I was holding for that old successor")
		c.releaseSuccessorReplica()
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	curDate := time.Now()
	c.Monitor.AddContact(candidate, curDate)
	c.logger.WriteToFileOK("Calling SetSuccessor Method, setting %v at date %v", candidate, curDate)
	c.successor = ContactWithData[T]{Contact: candidate, Data: make(Store)}
}

func (c *BruteChord[T]) GetSuccessorSuccessor() T {
	c.lock.Lock()
	defer c.lock.Unlock()
	successorSuccessor := c.successorSuccessor.Contact
	return successorSuccessor
}

func (c *BruteChord[T]) GetPredecessor() T {
	c.lock.Lock()
	defer c.lock.Unlock()
	pred := c.predecessorRef
	return pred
}

func (c *BruteChord[T]) SetSuccessorSuccessor(candidate T) {
	if c.GetSuccessorSuccessor().getNodeId() != candidate.getNodeId() {
		c.logger.WriteToFileOK("Because I'll have a new successorSuccessor, I'll have to release the replicas I was holding for that old successorSuccessor")
		c.releaseSuccessorSuccessorReplica()
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	curDate := time.Now()
	c.Monitor.AddContact(candidate, curDate)
	c.logger.WriteToFileOK("Calling SetSuccessorSuccessor Method, setting %v at date %v", candidate, curDate)
	c.successorSuccessor = ContactWithData[T]{Contact: candidate, Data: make(Store)}
}

func (c *BruteChord[T]) SetPredecessor(candidate T) {
	c.lock.Lock()
	defer c.lock.Unlock()
	curDate := time.Now()
	c.Monitor.AddContact(candidate, curDate)
	c.logger.WriteToFileOK("Calling SetPredecessor Method, setting %v at date %v", candidate, curDate)
	c.predecessorRef = candidate
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
			go c.StabilizeStore()       // Stabilize My Data, because due to changes in successor I may have data that is not mine.
		}
	}
}
