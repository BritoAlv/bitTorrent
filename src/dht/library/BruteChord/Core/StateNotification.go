package Core

type TellMeYourState[contact Contact] struct {
	QueryHost contact
}

func (t TellMeYourState[contact]) HandleNotification(b *BruteChord[contact]) {
	b.logger.WriteToFileOK("Handling TellMeYourState from %v", t.QueryHost.GetNodeId())
	b.clientChordCommunication.SendRequest(ClientTask[contact]{
		Targets: []contact{t.QueryHost},
		Data: TellMeYourStateResponse[contact]{
			Sender: b.GetContact(0),
			State:  b.GetState(),
		},
	})
}

type TellMeYourStateResponse[contact Contact] struct {
	Sender contact
	State  string
}

func (t TellMeYourStateResponse[contact]) HandleNotification(b *BruteChord[contact]) {
	return
}
