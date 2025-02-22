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

type ReceiveDataReplicate[contact Contact] struct {
	DataOwner      contact
	SuccessorIndex int   // The sender of this is my SuccessorIndex successor, my first successor, or my second successor.
	TaskId         int64 // The ID of the task so that it can be confirmed.
	Data           Store
}

type ConfirmReplication[contact Contact] struct {
	TaskId int64
}

func (c ConfirmReplication[contact]) HandleNotification(b *BruteChord[contact]) {
	b.logger.WriteToFileOK("Handling ConfirmReplication with TaskId %v", c.TaskId)
	b.logger.WriteToFileOK("I will now delete the data from my local storage")
	b.SetPendingResponse(c.TaskId, Confirmations{Confirmation: true, Value: nil})
}

func (r ReceiveDataReplicate[contact]) HandleNotification(b *BruteChord[contact]) {
	b.logger.WriteToFileOK("Handling ReceiveDataReplicate from %v, supposed to be my %v successor, taskId  %v", r.DataOwner.getNodeId(), r.SuccessorIndex, r.TaskId)
	if r.SuccessorIndex == 1 {
		if b.successor.Contact.getNodeId() == r.DataOwner.getNodeId() {
			b.logger.WriteToFileOK("I am the first successor of %v", r.DataOwner.getNodeId())
			b.logger.WriteToFileOK("I will now store the data")
			b.ReplaceStore(&b.successor.Data, r.Data)
			b.ClientChordCommunication.sendRequest(ClientTask[contact]{
				Targets: []contact{r.DataOwner},
				Data:    ConfirmReplication[contact]{TaskId: r.TaskId},
			})
		} else {
			b.logger.WriteToFileOK("I am not the first successor of %v", r.DataOwner.getNodeId())
		}
	} else if r.SuccessorIndex == 2 {
		if b.successorSuccessor.Contact.getNodeId() == r.DataOwner.getNodeId() {
			b.logger.WriteToFileOK("I am the second successor of %v", r.DataOwner.getNodeId())
			b.logger.WriteToFileOK("I will now store the data")
			b.ReplaceStore(&b.successorSuccessor.Data, r.Data)
			b.ClientChordCommunication.sendRequest(ClientTask[contact]{
				Targets: []contact{r.DataOwner},
				Data:    ConfirmReplication[contact]{TaskId: r.TaskId},
			})
		} else {
			b.logger.WriteToFileOK("I am not the second successor of %v", r.DataOwner.getNodeId())
		}
	}
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
		if b.GetSuccessor().getNodeId() == a.Contact.getNodeId() {
			b.logger.WriteToFileOK("I am already its predecessor of %v", a.Contact.getNodeId())
		} else {
			b.SetSuccessor(a.Contact)
			b.SetSuccessorSuccessor(a.MySuccessor)
			b.logger.WriteToFileOK("I am now the predecessor of %v,  so I will tell him that Im its predecessor and my predecessor is %v", a.Contact.getNodeId(), b.predecessorRef.getNodeId())
			b.logger.WriteToFileOK("My successor is now %v and my second successor is %v", b.GetSuccessor().getNodeId(), b.GetSuccessorSuccessor().getNodeId())
			b.ClientChordCommunication.sendRequest(ClientTask[contact]{
				Targets: []contact{a.Contact},
				Data: ImYourPredecessor[contact]{
					MyContact:     b.GetContact(),
					MyPredecessor: b.GetPredecessor(),
				},
			})
		}
	} else {
		b.logger.WriteToFileOK("I am not the predecessor of %v", a.Contact.getNodeId())
	}
}
