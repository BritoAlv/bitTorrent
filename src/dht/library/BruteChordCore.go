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
	MyData                        Store                   // My Data That I'm actually Responsible for.
	PendingResponses              map[int64]Confirmations // I need to check if someone answered the request I sent.
}

func NewBruteChord[T Contact](serverChordCommunication Server[T], clientChordCommunication Client[T], monitor Monitor[T]) *BruteChord[T] {
	var node = BruteChord[T]{}
	node.id = GenerateRandomBinaryId()
	node.lock = sync.Mutex{}
	node.DeadContacts = make([]T, 0)
	node.logger = *common.NewLogger(ToString(node.id) + ".txt")
	node.NotificationChannelServerNode = make(chan Notification[T])
	node.ServerChordCommunication = serverChordCommunication
	node.ServerChordCommunication.SetData(node.NotificationChannelServerNode, node.id)
	node.ClientChordCommunication = clientChordCommunication
	node.Monitor = monitor
	node.SetSuccessor(node.DefaultSuccessor())
	node.SetPredecessor(node.DefaultSuccessor())
	node.PendingResponses = make(map[int64]Confirmations)
	return &node
}

func (c *BruteChord[T]) GetId() ChordHash {
	return c.id
}

func (c *BruteChord[T]) GetData(key ChordHash) []byte {
	c.lock.Lock()
	data := c.MyData[key]
	c.lock.Unlock()
	return data
}

func (c *BruteChord[T]) GetAllData() Store {
	c.lock.Lock()
	data := c.MyData // this is a copy or a reference ?
	c.lock.Unlock()
	return data
}

func (c *BruteChord[T]) SetData(key ChordHash, value []byte) {
	c.lock.Lock()
	c.MyData[key] = value
	c.lock.Unlock()
}

func (c *BruteChord[T]) ReplaceStore(placeHolder *Store, data Store) {
	c.lock.Lock()
	*placeHolder = data
	c.lock.Unlock()
}

func (c *BruteChord[T]) GetSuccessor() T {
	c.lock.Lock()
	c.logger.WriteToFileOK("Calling GetSuccessor Method, returning %v", c.successor)
	successor := c.successor.Contact
	c.lock.Unlock()
	return successor
}

func (c *BruteChord[T]) SetPendingResponse(taskId int64, confirmation Confirmations) {
	c.lock.Lock()
	c.PendingResponses[taskId] = confirmation
	c.lock.Unlock()
}

func (c *BruteChord[T]) GetPendingResponse(taskId int64) Confirmations {
	c.lock.Lock()
	confirmation, _ := c.PendingResponses[taskId]
	c.lock.Unlock()
	return confirmation
}

func (c *BruteChord[T]) releaseReplicas(placeHolder *Store) {
	c.lock.Lock()
	for key, data := range *placeHolder {
		c.MyData[key] = data
		delete(*placeHolder, key)
	}
	c.lock.Unlock()
}

func (c *BruteChord[T]) SetSuccessor(candidate T) {
	if c.successor.Contact.getNodeId() != candidate.getNodeId() {
		c.logger.WriteToFileOK("Because I'll have a new successor, I'll have to release the replicas I was holding for that old successor")
		c.releaseReplicas(&c.successor.Data)
	}
	c.lock.Lock()
	curDate := time.Now()
	c.Monitor.AddContact(candidate, curDate)
	c.logger.WriteToFileOK("Calling SetSuccessor Method, setting %v at date %v", candidate, curDate)
	c.successor = ContactWithData[T]{Contact: candidate, Data: make(Store)}
	c.lock.Unlock()
}

func (c *BruteChord[T]) GetSuccessorSuccessor() T {
	c.lock.Lock()
	successorSuccessor := c.successorSuccessor.Contact
	c.lock.Unlock()
	return successorSuccessor
}

func (c *BruteChord[T]) GetPredecessor() T {
	c.lock.Lock()
	pred := c.predecessorRef
	c.lock.Unlock()
	return pred
}

func (c *BruteChord[T]) SetSuccessorSuccessor(candidate T) {
	if c.successorSuccessor.Contact.getNodeId() != candidate.getNodeId() {
		c.logger.WriteToFileOK("Because I'll have a new successorSuccessor, I'll have to release the replicas I was holding for that old successorSuccessor")
		c.releaseReplicas(&c.successorSuccessor.Data)
	}
	c.lock.Lock()
	curDate := time.Now()
	c.Monitor.AddContact(candidate, curDate)
	c.logger.WriteToFileOK("Calling SetSuccessorSuccessor Method, setting %v at date %v", candidate, curDate)
	c.successorSuccessor = ContactWithData[T]{Contact: candidate, Data: make(Store)}
	c.lock.Unlock()
}

func (c *BruteChord[T]) SetPredecessor(candidate T) {
	c.lock.Lock()
	curDate := time.Now()
	c.Monitor.AddContact(candidate, curDate)
	c.logger.WriteToFileOK("Calling SetPredecessor Method, setting %v at date %v", candidate, curDate)
	c.predecessorRef = candidate
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
			go c.StabilizeStore()       // Stabilize My Data, because due to changes in successor I may have data that is not mine.
		}
	}
}
