package Core

type getRequest[contact Contact] struct {
	QueryHost contact
	GetId     int64
	Key       ChordHash
}

type receivedGetRequest[contact Contact] struct {
	Sender contact
	GetId  int64
	Value  []byte
}

func (g getRequest[contact]) HandleNotification(b *BruteChord[contact]) {
	b.logger.WriteToFileOK("Handling getRequest from %v with GetId = %v", g.QueryHost.GetNodeId(), g.GetId)
	if b.responsible(g.Key) {
		clientTask := ClientTask[contact]{
			Targets: []contact{g.QueryHost},
			Data: receivedGetRequest[contact]{
				Sender: g.QueryHost,
				GetId:  g.GetId,
				Value:  b.getData(g.Key, 0),
			},
		}
		b.logger.WriteToFileOK("Sending confirmations to %v with GetId = %v", g.QueryHost.GetNodeId(), g.GetId)
		b.clientChordCommunication.SendRequest(clientTask)
	} else {
		clientTask := ClientTask[contact]{
			Targets: []contact{b.GetContact(1)},
			Data:    g,
		}
		b.clientChordCommunication.SendRequest(clientTask)
	}
}

func (r receivedGetRequest[contact]) HandleNotification(b *BruteChord[contact]) {
	b.logger.WriteToFileOK("Handling receivedGetRequest from %v with GetId = %v", r.Sender.GetNodeId(), r.GetId)
	b.setPendingResponse(r.GetId, confirmations{
		Confirmation: true,
		Value:        r.Value,
	})
}
