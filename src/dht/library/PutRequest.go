package library

/*
Implementation for Put Notifications is the following:

- if node is responsible for the key, then it will store the value and send a ReceivedPutRequest to the QueryHost, in other
case it will forward the request to its successor.
- ReceivedPutRequest will set the PendingConfirmation for the PutId to True in the QueryHost node, that node should check
that periodically to confirm the result of the query.
*/

type PutRequest[contact Contact] struct {
	QueryHost contact
	PutId     int64
	Key       ChordHash
	Value     []byte
}

type ReceivedPutRequest[contact Contact] struct {
	Sender contact
	PutId  int64
}

func (r ReceivedPutRequest[contact]) HandleNotification(b *BruteChord[contact]) {
	b.logger.WriteToFileOK("Handling ReceivedPutRequest from %v with PutId = %v", r.Sender.getNodeId(), r.PutId)
	b.lock.Lock()
	b.logger.WriteToFileOK("Setting PendingConfirmation for PutId = %v to True", r.PutId)
	b.PendingResponses[r.PutId] = Confirmations{
		Confirmation: true,
		Value:        nil,
	}
	b.lock.Unlock()
}

func (p PutRequest[contact]) HandleNotification(b *BruteChord[contact]) {
	// send a Received put request
	b.logger.WriteToFileOK("Handling PutRequest from %v with PutId = %v", p.QueryHost.getNodeId(), p.PutId)
	b.lock.Lock()
	between := Between(b.GetId(), p.Key, b.GetSuccessor().getNodeId())
	b.lock.Unlock()
	if between {
		b.logger.WriteToFileOK("Sending ReceivedPutRequest to %v with PutId = %v", p.QueryHost.getNodeId(), p.PutId)
		b.lock.Lock()
		b.MyData[p.Key] = p.Value
		b.lock.Unlock()
		b.ClientChordCommunication.sendRequest(ClientTask[contact]{
			Targets: []contact{p.QueryHost},
			Data: ReceivedPutRequest[contact]{
				Sender: p.QueryHost,
				PutId:  p.PutId,
			},
		})
	} else {
		b.ClientChordCommunication.sendRequest(ClientTask[contact]{
			Targets: []contact{b.GetSuccessor()},
			Data:    p,
		})
	}
}
