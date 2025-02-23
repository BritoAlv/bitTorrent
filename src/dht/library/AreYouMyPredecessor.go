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
	b.SetPredecessor(i.MyContact)
	myPredecessor := i.MyContact
	myPredecessorPredecessor := i.MyPredecessor
	b.logger.WriteToFileOK("My predecessor is now %v and its predecessor is %v", myPredecessor.getNodeId(), myPredecessorPredecessor.getNodeId())
	b.logger.WriteToFileOK("I will now replicate my data to my new predecessors")

	b.ReplicateData(1, myPredecessor)
	b.ReplicateData(2, myPredecessorPredecessor)
}



func (a AreYouMyPredecessor[contact]) HandleNotification(b *BruteChord[contact]) {
	// update my settings only if needed.
	b.logger.WriteToFileOK("Handling AreYouMyPredecessor from %v", a.Contact.getNodeId())
	b.logger.WriteToFileOK("My ID is %v = %v, Query comes from Node with ID %v = %v, and my successor ID is %v = %v", b.GetContact().getNodeId(), b.GetContact().getNodeId(), a.Contact.getNodeId(), a.Contact.getNodeId(), b.GetSuccessor().getNodeId(), b.GetSuccessor().getNodeId())
	if a.Contact.getNodeId() == b.GetId() {
		b.logger.WriteToFileOK("Ignoring the request because I am the sender")
		return
	}
	// If the Node asking is between me and my successor, then I am the predecessor of that Node.
	if Between(b.GetId(), a.Contact.getNodeId(), b.GetSuccessor().getNodeId()) {
		b.logger.WriteToFileOK("I am the predecessor of %v", a.Contact.getNodeId())
		if b.GetSuccessor().getNodeId() == a.Contact.getNodeId() && b.GetSuccessorSuccessor().getNodeId() == a.MySuccessor.getNodeId() {
			b.logger.WriteToFileOK("I am already its predecessor of %v, and I have also already its predecesor so nothing new", a.Contact.getNodeId())
		} else {
			b.SetSuccessor(a.Contact)
			b.SetSuccessorSuccessor(a.MySuccessor)
			b.logger.WriteToFileOK("I am now the predecessor of %v,  so I will tell him that Im its predecessor and my predecessor is %v", a.Contact.getNodeId(), b.GetPredecessor().getNodeId())
			b.logger.WriteToFileOK("My successor is now %v and my second successor is %v", b.GetSuccessor().getNodeId(), b.GetSuccessorSuccessor().getNodeId())
		}
		// Confirm in any case it is its predecessor.
		b.ClientChordCommunication.sendRequest(ClientTask[contact]{
			Targets: []contact{a.Contact},
			Data: ImYourPredecessor[contact]{
				MyContact:     b.GetContact(),
				MyPredecessor: b.GetPredecessor(),
			},
		})
	} else {
		b.logger.WriteToFileOK("I am not the predecessor of %v", a.Contact.getNodeId())
	}
}