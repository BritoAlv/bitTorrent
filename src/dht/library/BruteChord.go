package library

import (
	"bittorrent/common"
	"math/rand/v2"
	"reflect"
	"strconv"
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

func generateTaskId() int64 {
	return rand.Int64()
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
	node.logger = *common.NewLogger(strconv.Itoa(BinaryArrayToInt(node.id)) + ".txt")
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

func (c *BruteChord[T]) GetSuccessor() T {
	c.lock.Lock()
	c.logger.WriteToFileOK("Calling GetSuccessor Method, returning %v", c.successor)
	successor := c.successor.Contact
	c.lock.Unlock()
	return successor
}

func (c *BruteChord[T]) SendRequestUntilConfirmation(clientTask ClientTask[T], taskId int64) {
	c.logger.WriteToFileOK("Calling SendRequestUntilConfirmation Method with request %v and taskId %v", clientTask, taskId)

	c.lock.Lock()
	c.PendingResponses[taskId] = Confirmations{Confirmation: false, Value: nil}
	c.lock.Unlock()
	for i := 0; i < Attempts; i++ {
		c.lock.Lock()
		val, _ := c.PendingResponses[taskId]
		c.lock.Unlock()
		if val.Confirmation {
			c.logger.WriteToFileOK("Received Confirmation for taskId %v", taskId)
			return
		}
		c.ClientChordCommunication.sendRequest(clientTask)
		time.Sleep(WaitingTime * time.Second)
	}
	c.logger.WriteToFileOK("Didn't receive Confirmation for taskId %v", taskId)
}

func (c *BruteChord[T]) ReplicateData(successorIndex int, target T) {
	c.lock.Lock()
	c.logger.WriteToFileOK("Calling ReplicateData Method with successorIndex %v designated to %v", successorIndex, target)

	taskId := generateTaskId()
	clientTask := ClientTask[T]{
		Targets: []T{target},
		Data: ReceiveDataReplicate[T]{
			SuccessorIndex: successorIndex,
			TaskId:         taskId,
			Data:           c.MyData,
			DataOwner:      c.GetContact(),
		},
	}
	go c.SendRequestUntilConfirmation(clientTask, taskId)
}

func (c *BruteChord[T]) SetSuccessor(candidate T) {
	c.lock.Lock()
	curDate := time.Now()
	c.Monitor.AddContact(candidate, curDate)
	c.logger.WriteToFileOK("Calling SetSuccessor Method, setting %v at date %v", candidate, curDate)
	if c.successor.Contact.getNodeId() != candidate.getNodeId() {
		c.logger.WriteToFileOK("Because I'll have a new successor, I'll have to release the replicas I was holding for that old successor")
		for key, data := range c.successor.Data {
			c.MyData[key] = data
		}
	}
	c.successor = ContactWithData[T]{Contact: candidate, Data: make(Store)}
	c.lock.Unlock()
}

func (c *BruteChord[T]) Get(key ChordHash) []byte {
	c.logger.WriteToFileOK("Calling Get Method on key %v", BinaryArrayToInt(key))
	taskId := generateTaskId()
	clientTask := ClientTask[T]{
		Targets: []T{c.GetSuccessor()},
		Data: GetRequest[T]{
			QueryHost: c.GetContact(),
			GetId:     taskId,
			Key:       key,
		},
	}
	c.SendRequestUntilConfirmation(clientTask, taskId)
	return c.PendingResponses[taskId].Value
}

func (c *BruteChord[T]) Put(key ChordHash, value []byte) bool {
	c.logger.WriteToFileOK("Calling Put Method with key %v", BinaryArrayToInt(key))
	// create the taskId waiting for the response and send a put request to yourself. It is the same logic for the Get now.
	taskId := generateTaskId()
	taskClient := ClientTask[T]{
		Targets: []T{c.GetSuccessor()},
		Data: PutRequest[T]{
			QueryHost: c.GetContact(),
			PutId:     taskId,
			Key:       key,
			Value:     value,
		},
	}
	c.SendRequestUntilConfirmation(taskClient, taskId)
	return c.PendingResponses[taskId].Confirmation
}

func (c *BruteChord[T]) StabilizeStore() {
	curDate := time.Now()
	c.lock.Lock()
	c.logger.WriteToFileOK("Stabilizing Store at time %v", curDate)
	for key, data := range c.MyData {
		c.logger.WriteToFileOK("Calling Put on key %v", key)
		go c.Put(key, data)
		delete(c.MyData, key)
		c.logger.WriteToFileOK("Deleted key %v from Store", key)
	}
	c.lock.Unlock()
}

func (c *BruteChord[T]) GetSuccessorSuccessor() T {
	return c.successorSuccessor.Contact
}

func (c *BruteChord[T]) GetPredecessor() T {
	c.lock.Lock()
	pred := c.predecessorRef
	c.lock.Unlock()
	return pred
}

func (c *BruteChord[T]) SetSuccessorSuccessor(candidate T) {
	curDate := time.Now()
	c.lock.Lock()
	c.Monitor.AddContact(candidate, curDate)
	c.logger.WriteToFileOK("Calling SetSuccessorSuccessor Method, setting %v at date %v", candidate, curDate)
	if c.successorSuccessor.Contact.getNodeId() != candidate.getNodeId() {
		for key, data := range c.successorSuccessor.Data {
			c.MyData[key] = data
		}
	}
	c.successorSuccessor = ContactWithData[T]{Contact: candidate, Data: make(Store)}
	c.lock.Unlock()
}

func (c *BruteChord[T]) SetPredecessor(candidate T) {
	curDate := time.Now()
	c.lock.Lock()
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

func (c *BruteChord[T]) sendCheckPredecessor() {
	c.logger.WriteToFileOK("Calling sendCheckPredecessor Method")
	c.logger.WriteToFileOK("Sending AreYouMyPredecessor Notification to Everyone")
	c.ClientChordCommunication.sendRequestEveryone(AreYouMyPredecessor[T]{
		Contact:     c.GetContact(),
		MySuccessor: c.GetSuccessorSuccessor(),
	})
}

func (c *BruteChord[T]) sendCheckAlive() {
	c.logger.WriteToFileOK("Calling sendCheckAlive Method")
	c.logger.WriteToFileOK("Sending AreYouAliveNotification to %v", c.GetSuccessor())
	c.ClientChordCommunication.sendRequest(ClientTask[T]{
		Targets: []T{c.GetSuccessor(), c.successorSuccessor.Contact},
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
	successorSuccessor := c.GetSuccessorSuccessor()
	if !c.Monitor.CheckAlive(successorSuccessor, 3*WaitingTime) {
		c.logger.WriteToFileOK("My Successor Successor %v looks Dead to me", successorSuccessor)
		c.DeadContacts = append(c.DeadContacts, successorSuccessor)
		c.Monitor.DeleteContact(successorSuccessor)
		c.SetSuccessorSuccessor(c.GetSuccessor())
	}
}

func (c *BruteChord[T]) StopWorking() {
	c.NotificationChannelServerNode <- nil
}
