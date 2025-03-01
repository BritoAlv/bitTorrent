package Core

import (
	"bittorrent/common"
	"reflect"
	"sync"
	"time"
)

type BruteChord[T Contact] struct {
	id                            ChordHash // The ID of the Node.
	info                          [3]contactWithData[T]
	predecessorRef                T                       // Keep Track of the Predecessor, due to replication.
	pendingResponses              map[int64]confirmations // I need to check if someone answered the request I sent.
	lock                          sync.Mutex              // The Pointers shouldn't be updated concurrently.
	monitor                       Monitor[T]              // To Keep Track of HeartBeats.
	notificationChannelServerNode chan Notification[T]    // A channel that will be intermediary between the Server and the Node.
	serverChordCommunication      Server[T]               // A Server that will receive notifications from contacts of type T.
	clientChordCommunication      Client[T]               // A Client that will send notifications to others nodes of type T.
	logger                        common.Logger           // To Log Everything The Node is doing.
	isWorking                     bool
}

func NewBruteChord[T Contact](serverChordCommunication Server[T], clientChordCommunication Client[T], monitor Monitor[T], id ChordHash) *BruteChord[T] {
	var node = BruteChord[T]{}
	node.id = id
	node.lock = sync.Mutex{}
	node.logger = *common.NewLogger(toString(node.id) + ".txt")
	node.notificationChannelServerNode = make(chan Notification[T])
	node.serverChordCommunication = serverChordCommunication
	node.serverChordCommunication.SetData(node.notificationChannelServerNode, node.id)
	node.clientChordCommunication = clientChordCommunication
	node.monitor = monitor
	node.pendingResponses = make(map[int64]confirmations)
	node.info[0].Contact = serverChordCommunication.GetContact()
	for i := 0; i < 3; i++ {
		node.info[i].Data = make(Store)
	}
	node.setContact(node.defaultSuccessor(), -1)
	node.setContact(node.defaultSuccessor(), 1)
	node.setContact(node.defaultSuccessor(), 2)
	node.SetWork(true)
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

func (c *BruteChord[T]) stabilizeOwnData() {
	c.lock.Lock()
	defer c.lock.Unlock()
	successor := c.info[1].Contact.GetNodeId()
	for key, value := range c.info[0].Data {
		if !between(c.id, key, successor) {
			delete(c.info[0].Data, key)
			go c.Put(key, value)
		}
	}
}

func (c *BruteChord[T]) getData(key ChordHash, index int) []byte {
	c.lock.Lock()
	defer c.lock.Unlock()
	data := c.copyValue(c.info[index].Data[key])
	return data
}

func (c *BruteChord[T]) getAllData(index int) Store {
	c.lock.Lock()
	defer c.lock.Unlock()
	data := c.copyStore(c.info[index].Data)
	return data
}

func (c *BruteChord[T]) setData(key ChordHash, value []byte, index int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.info[index].Data[key] = value
}

func (c *BruteChord[T]) setPendingResponse(taskId int64, confirmation confirmations) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.pendingResponses[taskId] = confirmation
}

func (c *BruteChord[T]) addNewData(data Store, index int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.info[index].Data = data
}

func (c *BruteChord[T]) getPendingResponse(taskId int64) confirmations {
	c.lock.Lock()
	defer c.lock.Unlock()
	confirmation, _ := c.pendingResponses[taskId]
	return confirmation
}

func (c *BruteChord[T]) setContact(candidate T, index int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if index == -1 {
		c.predecessorRef = candidate
	} else {
		if c.info[index].Contact.GetNodeId() != candidate.GetNodeId() {
			c.logger.WriteToFileOK("Because I'll have a new successor, I'll have to release the replicas I was holding for that old successor")
			for key, data := range c.info[index].Data {
				go c.Put(key, data)
			}
			c.info[index] = contactWithData[T]{Contact: candidate, Data: make(Store)}
		}
	}
	curDate := time.Now()
	c.monitor.AddContact(candidate, curDate)
}

func (c *BruteChord[T]) defaultSuccessor() T {
	contact := c.info[0].Contact
	c.logger.WriteToFileOK("Calling defaultSuccessor Method, returning %v", contact)
	return contact
}

func (c *BruteChord[T]) GetContact(index int) T {
	c.lock.Lock()
	defer c.lock.Unlock()
	var contact T
	if index == -1 {
		contact = c.predecessorRef
	} else {
		contact = c.info[index].Contact
	}
	return contact
}

func (c *BruteChord[T]) SetWork(value bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if !c.isWorking && value {
		c.logger.WriteToFileOK("Calling beginWorking Method at data %v", time.Now())
		go c.cpu()
	}
	if c.isWorking && !value {
		c.notificationChannelServerNode <- nil
	}
	c.isWorking = value
}

func (c *BruteChord[T]) cpu() {
	ticker := time.NewTicker(WaitingTime * time.Second)
	defer ticker.Stop() // Ensure the ticker stops when function exits.
	for {
		select {
		case notification := <-c.notificationChannelServerNode:
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
			go c.stabilizeOwnData()
		}
	}
}
