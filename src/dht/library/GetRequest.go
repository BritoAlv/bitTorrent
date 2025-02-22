package library

type GetRequest[contact Contact] struct {
	QueryHost contact
	GetId     int64
	Key       ChordHash
}

type ReceivedGetRequest[contact Contact] struct {
	Sender contact
	GetId  int64
	Value  []byte
}

func (g GetRequest[contact]) HandleNotification(b *BruteChord[contact]) {
	b.logger.WriteToFileOK("Handling GetRequest from %v with GetId = %v", g.QueryHost.getNodeId(), g.GetId)
	between := Between(b.GetId(), g.Key, b.GetSuccessor().getNodeId())
	if between {
		clientTask := ClientTask[contact]{
			Targets: []contact{g.QueryHost},
			Data: ReceivedGetRequest[contact]{
				Sender: g.QueryHost,
				GetId:  g.GetId,
				Value:  b.Get(g.Key),
			},
		}
		b.logger.WriteToFileOK("Sending Confirmations to %v with GetId = %v", g.QueryHost.getNodeId(), g.GetId)
		b.ClientChordCommunication.sendRequest(clientTask)
	} else {
		clientTask := ClientTask[contact]{
			Targets: []contact{b.GetSuccessor()},
			Data:    g,
		}
		b.ClientChordCommunication.sendRequest(clientTask)
	}
}

func (r ReceivedGetRequest[contact]) HandleNotification(b *BruteChord[contact]) {
	b.logger.WriteToFileOK("Handling ReceivedGetRequest from %v with GetId = %v", r.Sender.getNodeId(), r.GetId)
	b.SetPendingResponse(r.GetId, Confirmations{
		Confirmation: true,
		Value:        r.Value,
	})
}
