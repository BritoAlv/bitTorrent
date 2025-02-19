package library

import (
	"time"
)

type Notification[contact Contact] interface {
	HandleNotification(*BruteChord[contact])
}

type ImAliveNotification[contact Contact] struct {
	Contact contact
}

func (a ImAliveNotification[contact]) HandleNotification(b *BruteChord[contact]) {
	curDate := time.Now()
	b.logger.WriteToFileOK("Handling ImAliveNotification from %v at date %v", a.Contact.getNodeId(), curDate)
	b.logger.WriteToFileOK("Updating Contact Date of %v to %v", a.Contact.getNodeId(), curDate)
	b.Monitor.UpdateContactDate(a.Contact, curDate)
}

type AreYouAliveNotification[contact Contact] struct {
	Contact contact
}

func (a AreYouAliveNotification[contact]) HandleNotification(b *BruteChord[contact]) {
	b.logger.WriteToFileOK("Handling AreYouAliveNotification from %v", a.Contact.getNodeId())
	var target = make([]contact, 1)
	target[0] = a.Contact
	b.ClientChordCommunication.sendRequest(ClientTask[contact]{
		Targets: target,
		Data: ImAliveNotification[contact]{
			b.GetContact(),
		},
	})
}

type AreYouMyPredecessor[contact Contact] struct {
	Contact contact
}

func (a AreYouMyPredecessor[contact]) HandleNotification(b *BruteChord[contact]) {
	// update my settings only if needed.
	b.logger.WriteToFileOK("Handling AreYouMyPredecessor from %v", a.Contact.getNodeId())
	b.logger.WriteToFileOK("My ID is %v = %v, Query comes from Node with ID %v = %v, and my successorList ID is %v = %v", b.GetContact().getNodeId(), BinaryArrayToInt(b.GetContact().getNodeId()), a.Contact.getNodeId(), BinaryArrayToInt(a.Contact.getNodeId()), b.GetSuccessor().getNodeId(), BinaryArrayToInt(b.GetSuccessor().getNodeId()))
	if a.Contact.getNodeId() == b.GetId() {
		b.logger.WriteToFileOK("Ignoring the request because I am the sender")
		return
	}
	// If the Node asking is between me and my successorList, then I am the predecessor of that Node.
	if Between(b.GetId(), a.Contact.getNodeId(), b.GetSuccessor().getNodeId()) {
		b.logger.WriteToFileOK("I am the predecessor of %v", a.Contact.getNodeId())
		if b.GetSuccessor().getNodeId() == a.Contact.getNodeId() {
			b.logger.WriteToFileOK("I am already its predecessor of %v", a.Contact.getNodeId())
			return
		} else {
			b.SetSuccessor(a.Contact)
		}
	} else {
		b.logger.WriteToFileOK("I am not the predecessor of %v", a.Contact.getNodeId())
	}
}
