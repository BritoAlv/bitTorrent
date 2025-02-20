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
	b.lock.Lock()
	between := Between(b.GetId(), g.Key, b.GetSuccessor().getNodeId())
	b.lock.Unlock()
	if between {
		b.logger.WriteToFileOK("Sending Confirmations to %v with GetId = %v", g.QueryHost.getNodeId(), g.GetId)
		b.ClientChordCommunication.sendRequest(ClientTask[contact]{
			Targets: []contact{g.QueryHost},
			Data: ReceivedGetRequest[contact]{
				Sender: g.QueryHost,
				GetId:  g.GetId,
				Value:  b.Get(g.Key),
			},
		})
	} else {
		b.ClientChordCommunication.sendRequest(ClientTask[contact]{
			Targets: []contact{b.GetSuccessor()},
			Data:    g,
		})
	}
}

func (r ReceivedGetRequest[contact]) HandleNotification(b *BruteChord[contact]) {
	b.logger.WriteToFileOK("Handling ReceivedGetRequest from %v with GetId = %v", r.Sender.getNodeId(), r.GetId)
	b.lock.Lock()
	b.PendingResponses[r.GetId] = Confirmations{
		Confirmation: true,
		Value:        r.Value,
	}
	b.lock.Unlock()
}
