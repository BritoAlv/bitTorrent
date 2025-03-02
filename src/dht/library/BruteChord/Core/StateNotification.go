package Core

import (
	"fmt"
	"strconv"
)

type TellMeYourState[contact Contact] struct {
	QueryHost contact
}

type NodeState[contact Contact] struct {
	NodeId                 ChordHash
	SuccessorId            ChordHash
	SuccessorData          Store
	SuccessorSuccessorId   ChordHash
	SuccessorSuccessorData Store
	PredecessorId          ChordHash
	OwnData                Store
}

func (n NodeState[contact]) String() string {
	state := "Node: " + strconv.Itoa(int(n.NodeId)) + "\n"
	state += "Successor: " + strconv.Itoa(int(n.SuccessorId)) + "\n"
	state += "Successor Data Replicas Are: " + "\n"
	for _, key := range SortKeys(n.SuccessorData) {
		state += strconv.Itoa(int(key)) + " -> " + fmt.Sprintf("%v", n.SuccessorData[key]) + "\n"
	}
	state += "SuccessorSuccessor: " + strconv.Itoa(int(n.SuccessorSuccessorId)) + "\n"
	state += "SuccessorSuccessor Data Replica:" + "\n"
	for _, key := range SortKeys(n.SuccessorSuccessorData) {
		state += strconv.Itoa(int(key)) + " -> " + fmt.Sprintf("%v", n.SuccessorSuccessorData[key]) + "\n"
	}
	state += "Predecessor: " + strconv.Itoa(int(n.PredecessorId)) + "\n"
	state += "Data stored:\n"
	for _, key := range SortKeys(n.OwnData) {
		state += strconv.Itoa(int(key)) + " -> " + fmt.Sprintf("%v", n.OwnData[key]) + "\n"
	}
	return state
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
	State  NodeState[contact]
}

func (t TellMeYourStateResponse[contact]) HandleNotification(b *BruteChord[contact]) {
	return
}
