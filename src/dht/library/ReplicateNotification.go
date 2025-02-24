package library

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
		bSuccessor := b.GetContact(1)
		if bSuccessor.getNodeId() == r.DataOwner.getNodeId() {
			b.logger.WriteToFileOK("I am the first successor of %v", r.DataOwner.getNodeId())
			b.logger.WriteToFileOK("I will now store the data")
			b.AddNewData(r.Data, 1)
			b.ClientChordCommunication.sendRequest(ClientTask[contact]{
				Targets: []contact{r.DataOwner},
				Data:    ConfirmReplication[contact]{TaskId: r.TaskId},
			})
		} else {
			b.logger.WriteToFileOK("I am not the first successor of %v", r.DataOwner.getNodeId())
		}
	} else if r.SuccessorIndex == 2 {
		bSuccessorSuccessor := b.GetContact(2)
		if bSuccessorSuccessor.getNodeId() == r.DataOwner.getNodeId() {
			b.logger.WriteToFileOK("I am the second successor of %v", r.DataOwner.getNodeId())
			b.logger.WriteToFileOK("I will now store the data")
			b.AddNewData(r.Data, 2)
			b.ClientChordCommunication.sendRequest(ClientTask[contact]{
				Targets: []contact{r.DataOwner},
				Data:    ConfirmReplication[contact]{TaskId: r.TaskId},
			})
		} else {
			b.logger.WriteToFileOK("I am not the second successor of %v", r.DataOwner.getNodeId())
		}
	}
}
