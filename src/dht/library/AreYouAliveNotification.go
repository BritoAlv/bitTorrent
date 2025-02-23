package library

import (
	"time"
)

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

type ImAliveNotification[contact Contact] struct {
	Contact contact
}

func (a ImAliveNotification[contact]) HandleNotification(b *BruteChord[contact]) {
	curDate := time.Now()
	b.logger.WriteToFileOK("Handling ImAliveNotification from %v at date %v", a.Contact.getNodeId(), curDate)
	b.logger.WriteToFileOK("Updating Contact Date of %v to %v", a.Contact.getNodeId(), curDate)
	b.Monitor.UpdateContactDate(a.Contact, curDate)
}
