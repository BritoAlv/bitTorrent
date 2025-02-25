package library

type AreYouMyPredecessor[contact Contact] struct {
	Contact     contact
	MySuccessor contact
}

type ImYourPredecessor[contact Contact] struct {
	MyContact     contact
	MyPredecessor contact
}

func (i ImYourPredecessor[contact]) HandleNotification(b *BruteChord[contact]) {
	b.logger.WriteToFileOK("Handling ImYourPredecessor from %v", i.MyContact.getNodeId())
	b.SetContact(i.MyContact, -1)
	myPredecessor := i.MyContact
	myPredecessorPredecessor := i.MyPredecessor
	b.logger.WriteToFileOK("My predecessor is now %v and its predecessor is %v", myPredecessor.getNodeId(), myPredecessorPredecessor.getNodeId())
	b.logger.WriteToFileOK("I will now replicate my data to my new predecessors")
	b.ReplicateData(1, myPredecessor)
	b.ReplicateData(2, myPredecessorPredecessor)
}

func (a AreYouMyPredecessor[contact]) HandleNotification(b *BruteChord[contact]) {
	// update my settings only if needed.
	if a.Contact.getNodeId() == b.GetId() {
		b.logger.WriteToFileOK("Ignoring the request because I am the sender")
		return
	}
	// If the Node asking is between me and my successor, then I am the predecessor of that Node.
	if b.Responsible(a.Contact.getNodeId()) {
		b.logger.WriteToFileOK("I am the predecessor of %v", a.Contact.getNodeId())
		if b.GetContact(1).getNodeId() == a.Contact.getNodeId() && b.GetContact(2).getNodeId() == a.MySuccessor.getNodeId() {
			b.logger.WriteToFileOK("I am already its predecessor of %v, and I have also already its predecesor so nothing new", a.Contact.getNodeId())
		} else {
			b.SetContact(a.Contact, 1)
			b.SetContact(a.MySuccessor, 2)
		}
		// Confirm in any case I'm its predecessor.
		b.ClientChordCommunication.sendRequest(ClientTask[contact]{
			Targets: []contact{a.Contact},
			Data: ImYourPredecessor[contact]{
				MyContact:     b.GetContact(0),
				MyPredecessor: b.GetContact(-1),
			},
		})
	} else {
		b.logger.WriteToFileOK("I am not the predecessor of %v", a.Contact.getNodeId())
	}
}
