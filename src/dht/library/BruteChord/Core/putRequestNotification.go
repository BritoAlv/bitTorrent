package Core

/*
Implementation for Put Notifications is the following:

- if node is responsible for the key, then it will store the value and send a receivedPutRequest to the QueryHost, in other
case it will forward the request to its successor.
- receivedPutRequest will set the PendingConfirmation for the PutId to True in the QueryHost node, that node should check
that periodically to confirm the result of the query.
*/

type putRequest[contact Contact] struct {
	QueryHost contact
	PutId     int64
	Key       ChordHash
	Value     []byte
}

type receivedPutRequest[contact Contact] struct {
	Sender contact
	PutId  int64
}

func (r receivedPutRequest[contact]) HandleNotification(b *BruteChord[contact]) {
	b.logger.WriteToFileOK("Handling receivedPutRequest from %v with PutId = %v", r.Sender.GetNodeId(), r.PutId)
	b.logger.WriteToFileOK("Setting PendingConfirmation for PutId = %v to True", r.PutId)
	b.setPendingResponse(r.PutId, confirmations{
		Confirmation: true,
		Value:        nil,
	})
}

func (p putRequest[contact]) HandleNotification(b *BruteChord[contact]) {
	// send a Received put request
	b.logger.WriteToFileOK("Handling putRequest from %v with PutId = %v", p.QueryHost.GetNodeId(), p.PutId)
	if b.responsible(p.Key) {
		b.logger.WriteToFileOK("Sending receivedPutRequest to %v with PutId = %v", p.QueryHost.GetNodeId(), p.PutId)
		b.setData(p.Key, p.Value, 0)
		b.clientChordCommunication.SendRequest(ClientTask[contact]{
			Targets: []contact{p.QueryHost},
			Data: receivedPutRequest[contact]{
				Sender: p.QueryHost,
				PutId:  p.PutId,
			},
		})
	} else {
		clientTask := ClientTask[contact]{
			Targets: []contact{b.GetContact(1)},
			Data:    p,
		}
		b.clientChordCommunication.SendRequest(clientTask)
	}
}
