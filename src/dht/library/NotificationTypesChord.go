package library

import "time"

type Notification[contact Contact] interface {
	HandleNotification(*BruteChord[contact])
}

type ImAliveNotification[contact Contact] struct {
	Contact contact
}

func (a ImAliveNotification[contact]) HandleNotification(b *BruteChord[contact]) {
	b.Monitor.UpdateContactDate(a.Contact, time.Now())
}

type AreYouAliveNotification[contact Contact] struct {
	Contact contact
}

func (a AreYouAliveNotification[contact]) HandleNotification(b *BruteChord[contact]) {
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
	if Between(b.Id, a.Contact.getNodeId(), b.GetSuccessor().getNodeId()) || b.Id == b.GetSuccessor().getNodeId() {
		b.SetSuccessor(a.Contact)
	}
}

type KillNotification[contact Contact] struct {
	Contact contact
}

func (k KillNotification[contact]) HandleNotification(b *BruteChord[contact]) {
	if b.GetSuccessor().getNodeId() == k.Contact.getNodeId() {
		b.SetSuccessor(b.DefaultSuccessor())
	}
}
