package library

import (
	"fmt"
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
	b.Logger.WriteToFileOK(fmt.Sprintf("Handling ImAliveNotification from %v at date %v", a.Contact.getNodeId(), curDate))
	b.Logger.WriteToFileOK(fmt.Sprintf("Updating Contact Date of %v to %v", a.Contact.getNodeId(), curDate))
	b.Monitor.UpdateContactDate(a.Contact, curDate)
}

type AreYouAliveNotification[contact Contact] struct {
	Contact contact
}

func (a AreYouAliveNotification[contact]) HandleNotification(b *BruteChord[contact]) {
	b.Logger.WriteToFileOK(fmt.Sprintf("Handling AreYouAliveNotification from %v", a.Contact.getNodeId()))
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
	b.Logger.WriteToFileOK(fmt.Sprintf("Handling AreYouMyPredecessor from %v", a.Contact.getNodeId()))
	b.Logger.WriteToFileOK(fmt.Sprintf("My ID is %v = %v, Query comes from Node with ID %v = %v, and my successor ID is %v = %v", b.GetContact().getNodeId(), BinaryArrayToInt(b.GetContact().getNodeId()), a.Contact.getNodeId(), BinaryArrayToInt(a.Contact.getNodeId()), b.GetSuccessor().getNodeId(), BinaryArrayToInt(b.GetSuccessor().getNodeId())))
	if a.Contact.getNodeId() == b.Id {
		b.Logger.WriteToFileOK(fmt.Sprintf("Ignoring the request because I am the sender"))
		return
	}
	// If the Node asking is between me and my successor, then I am the predecessor of that Node.
	if Between(b.Id, a.Contact.getNodeId(), b.GetSuccessor().getNodeId()) {
		b.Logger.WriteToFileOK(fmt.Sprintf("I am the predecessor of %v", a.Contact.getNodeId()))
		if b.GetSuccessor().getNodeId() == a.Contact.getNodeId() {
			b.Logger.WriteToFileOK(fmt.Sprintf("I am already its predecessor of %v", a.Contact.getNodeId()))
			return
		} else {
			b.SetSuccessor(a.Contact)
		}
	} else {
		b.Logger.WriteToFileOK(fmt.Sprintf("I am not the predecessor of %v", a.Contact.getNodeId()))
	}
}
