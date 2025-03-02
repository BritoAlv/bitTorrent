package Core

type areYouMyPredecessor[contact Contact] struct {
	Contact     contact
	MySuccessor contact
}

type imYourPredecessor[contact Contact] struct {
	MyContact     contact
	MyPredecessor contact
}

func (i imYourPredecessor[contact]) HandleNotification(b *BruteChord[contact]) {
	b.logger.WriteToFileOK("Handling imYourPredecessor from %v", i.MyContact.GetNodeId())
	b.setContact(i.MyContact, -1)
	myPredecessor := i.MyContact
	myPredecessorPredecessor := i.MyPredecessor
	b.logger.WriteToFileOK("My predecessor is now %v and its predecessor is %v", myPredecessor.GetNodeId(), myPredecessorPredecessor.GetNodeId())
	b.logger.WriteToFileOK("I will now replicate my data to my new predecessors")
	b.replicateData(1, myPredecessor)
	b.replicateData(2, myPredecessorPredecessor)
}

func (a areYouMyPredecessor[contact]) HandleNotification(b *BruteChord[contact]) {
	// update my settings only if needed.
	if a.Contact.GetNodeId() == b.GetId() {
		b.logger.WriteToFileOK("Ignoring the request because I am the sender")
		return
	}
	// If the Node asking is between me and my successor, then I am the predecessor of that Node.
	if b.responsible(a.Contact.GetNodeId()) {
		b.logger.WriteToFileOK("I am the predecessor of %v", a.Contact.GetNodeId())
		if b.GetContact(1).GetNodeId() == a.Contact.GetNodeId() && b.GetContact(2).GetNodeId() == a.MySuccessor.GetNodeId() {
			b.logger.WriteToFileOK("I am already its predecessor of %v, and I have also already its predecessor so nothing new", a.Contact.GetNodeId())
		} else {
			b.setContact(a.Contact, 1)
			b.setContact(a.MySuccessor, 2)
		}
		// Confirm in any case I'm its predecessor.
		b.clientChordCommunication.SendRequest(ClientTask[contact]{
			Targets: []contact{a.Contact},
			Data: imYourPredecessor[contact]{
				MyContact:     b.GetContact(0),
				MyPredecessor: b.GetContact(-1),
			},
		})
	} else {
		b.logger.WriteToFileOK("I am not the predecessor of %v", a.Contact.GetNodeId())
	}
}
