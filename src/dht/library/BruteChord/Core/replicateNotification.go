package Core

type receiveDataReplicate[contact Contact] struct {
	DataOwner      contact
	SuccessorIndex int   // The sender of this is my SuccessorIndex successor, my first successor, or my second successor.
	TaskId         int64 // The ID of the task so that it can be confirmed.
	Data           Store
}

type confirmReplication[contact Contact] struct {
	TaskId int64
}

func (c confirmReplication[contact]) HandleNotification(b *BruteChord[contact]) {
	b.logger.WriteToFileOK("Handling confirmReplication with TaskId %v", c.TaskId)
	b.logger.WriteToFileOK("I will now delete the data from my local storage")
	b.setPendingResponse(c.TaskId, confirmations{Confirmation: true, Value: nil})
}

func (r receiveDataReplicate[contact]) HandleNotification(b *BruteChord[contact]) {
	b.logger.WriteToFileOK("Handling receiveDataReplicate from %v, supposed to be my %v successor, taskId  %v", r.DataOwner.GetNodeId(), r.SuccessorIndex, r.TaskId)
	if r.SuccessorIndex == 1 {
		bSuccessor := b.GetContact(1)
		if bSuccessor.GetNodeId() == r.DataOwner.GetNodeId() {
			b.logger.WriteToFileOK("I am the first successor of %v", r.DataOwner.GetNodeId())
			b.logger.WriteToFileOK("I will now store the data")
			b.addNewData(r.Data, 1)
			b.clientChordCommunication.SendRequest(ClientTask[contact]{
				Targets: []contact{r.DataOwner},
				Data:    confirmReplication[contact]{TaskId: r.TaskId},
			})
		} else {
			b.logger.WriteToFileOK("I am not the first successor of %v", r.DataOwner.GetNodeId())
		}
	} else if r.SuccessorIndex == 2 {
		bSuccessorSuccessor := b.GetContact(2)
		if bSuccessorSuccessor.GetNodeId() == r.DataOwner.GetNodeId() {
			b.logger.WriteToFileOK("I am the second successor of %v", r.DataOwner.GetNodeId())
			b.logger.WriteToFileOK("I will now store the data")
			b.addNewData(r.Data, 2)
			b.clientChordCommunication.SendRequest(ClientTask[contact]{
				Targets: []contact{r.DataOwner},
				Data:    confirmReplication[contact]{TaskId: r.TaskId},
			})
		} else {
			b.logger.WriteToFileOK("I am not the second successor of %v", r.DataOwner.GetNodeId())
		}
	}
}
