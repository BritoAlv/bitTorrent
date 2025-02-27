package Core

import (
	"time"
)

type areYouAliveNotification[contact Contact] struct {
	Contact contact
}

func (a areYouAliveNotification[contact]) HandleNotification(b *BruteChord[contact]) {
	b.logger.WriteToFileOK("Handling areYouAliveNotification from %v", a.Contact.GetNodeId())
	var target = make([]contact, 1)
	target[0] = a.Contact
	b.clientChordCommunication.SendRequest(ClientTask[contact]{
		Targets: target,
		Data: imAliveNotification[contact]{
			b.GetContact(0),
		},
	})
}

type imAliveNotification[contact Contact] struct {
	Contact contact
}

func (a imAliveNotification[contact]) HandleNotification(b *BruteChord[contact]) {
	curDate := time.Now()
	b.logger.WriteToFileOK("Handling imAliveNotification from %v at date %v", a.Contact.GetNodeId(), curDate)
	b.logger.WriteToFileOK("Updating Contact Date of %v to %v", a.Contact.GetNodeId(), curDate)
	b.monitor.UpdateContactDate(a.Contact, curDate)
}
