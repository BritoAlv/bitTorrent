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
	b.Logger.WriteToFileOK(fmt.Sprintf("My ID is %v, Query comes from Node with ID %v, and my successor ID is %v", b.GetContact().getNodeId(), a.Contact.getNodeId(), b.GetSuccessor().getNodeId()))
	if Between(b.Id, a.Contact.getNodeId(), b.GetSuccessor().getNodeId()) || b.Id == b.GetSuccessor().getNodeId() {
		b.Logger.WriteToFileOK(fmt.Sprintf("I am the predecessor of %v", a.Contact.getNodeId()))
		b.SetSuccessor(a.Contact)
	} else {
		b.Logger.WriteToFileOK(fmt.Sprintf("I am not the predecessor of %v", a.Contact.getNodeId()))
	}
}
