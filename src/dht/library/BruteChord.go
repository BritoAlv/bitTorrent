package library

import "time"

type BruteChord[T Contact] struct {
	Id                            [NumberBits]uint8    // Every Node should have an ID.
	successor                     T                    // A Contact, hidden because Set and Get Methods will be overloaded.
	Monitor                       Monitor[T]           // To Keep Track of HeartBeats.
	NotificationChannelServerNode chan Notification[T] // A channel that will be intermediary between the Server and the Node.
	ServerChordCommunication      Server[T]            // A Server that will receive notifications from contacts of type T.
	ClientChordCommunication      Client[T]            // A Client that will send notifications to others nodes of type T.
}

func NewBruteChord[T Contact](serverChordCommunication Server[T], clientChordCommunication Client[T], monitor Monitor[T]) *BruteChord[T] {
	var node = BruteChord[T]{}
	node.Id = GenerateRandomBinaryId()
	node.NotificationChannelServerNode = make(chan Notification[T])
	node.ServerChordCommunication = serverChordCommunication
	node.ServerChordCommunication.SetData(node.NotificationChannelServerNode, node.Id)
	node.ClientChordCommunication = clientChordCommunication
	node.Monitor = monitor
	node.SetSuccessor(node.DefaultSuccessor())
	return &node
}

func (c *BruteChord[T]) GetSuccessor() T {

	return c.successor
}

func (c *BruteChord[T]) SetSuccessor(candidate T) {
	c.Monitor.AddContact(candidate, time.Now())
	c.successor = candidate
}

func (c *BruteChord[T]) DefaultSuccessor() T {
	return c.ServerChordCommunication.GetContact()
}

func (c *BruteChord[T]) GetContact() T {
	return c.ServerChordCommunication.GetContact()
}

// BeginWorking Callers should use a Barrier because this is an infinite loop.
func (c *BruteChord[T]) BeginWorking() {
	c.cpu()
}

func (c *BruteChord[T]) cpu() {
	for {
		select {
		case notification := <-c.NotificationChannelServerNode:
			notification.HandleNotification(c)
		case <-time.After(WaitingTime):
			go c.sendCheckPredecessor() // Stabilize Predecessor.
			go c.sendCheckAlive()       // Check if The Contacts I Have Are Alive.
			go c.killDead()             // Remove the Contacts that are Dead.
		}
	}
}

func (c *BruteChord[T]) sendCheckPredecessor() {
	c.ClientChordCommunication.sendRequestEveryone(AreYouMyPredecessor[T]{
		Contact: c.GetContact(),
	})
}

func (c *BruteChord[T]) sendCheckAlive() {
	c.ClientChordCommunication.sendRequest(ClientTask[T]{
		Targets: []T{c.GetSuccessor()},
		Data:    AreYouAliveNotification[T]{Contact: c.GetContact()},
	})
}

func (c *BruteChord[T]) killDead() {
	if !c.Monitor.CheckAlive(c.GetSuccessor(), 2*WaitingTime) {
		c.Monitor.DeleteContact(c.GetSuccessor())
		c.SetSuccessor(c.DefaultSuccessor())
	}
}
