package Core

type TellMeYourState[contact Contact] struct {
	QueryHost contact
}

func (t TellMeYourState[contact]) HandleNotification(b *BruteChord[contact]) {
	b.logger.WriteToFileOK("Handling TellMeYourState from %v", t.QueryHost.GetNodeId())
	b.clientChordCommunication.SendRequest(ClientTask[contact]{
		Targets: []contact{t.QueryHost},
		Data: TellMeYourStateResponse[contact]{
			State: b.GetState(),
		},
	})
}

type TellMeYourStateResponse[contact Contact] struct {
	State string
}

func (t TellMeYourStateResponse[contact]) HandleNotification(b *BruteChord[contact]) {
	return
}
