package library

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

func (c *BruteChord[T]) Responsible(key ChordHash) bool {
	return Between(c.GetId(), key, c.GetContact(1).getNodeId())
}

func (c *BruteChord[T]) SendRequestUntilConfirmation(clientTask ClientTask[T], taskId int64) {
	c.logger.WriteToFileOK("Calling SendRequestUntilConfirmation Method with request %v and taskId %v", clientTask, taskId)
	c.SetPendingResponse(taskId, Confirmations{Confirmation: false, Value: nil})
	for i := 0; i < Attempts; i++ {
		confirmation := c.GetPendingResponse(taskId)
		if confirmation.Confirmation {
			c.logger.WriteToFileOK("Received Confirmation for taskId %v", taskId)
			return
		}
		c.ClientChordCommunication.sendRequest(clientTask)
		time.Sleep(WaitingTime * time.Second)
	}
	c.logger.WriteToFileOK("Didn't receive Confirmation for taskId %v", taskId)
}

func (c *BruteChord[T]) ReplicateData(successorIndex int, target T) {
	c.logger.WriteToFileOK("Calling ReplicateData Method with successorIndex %v designated to %v", successorIndex, target)
	taskId := generateTaskId()
	clientTask := ClientTask[T]{
		Targets: []T{target},
		Data: ReceiveDataReplicate[T]{
			SuccessorIndex: successorIndex,
			TaskId:         taskId,
			Data:           c.GetAllData(0),
			DataOwner:      c.GetContact(0),
		},
	}
	go c.SendRequestUntilConfirmation(clientTask, taskId)
}

func (c *BruteChord[T]) Get(key ChordHash) ([]byte, bool) {
	c.logger.WriteToFileOK("Calling Get Method on key %v", key)
	taskId := generateTaskId()
	clientTask := ClientTask[T]{
		Targets: []T{c.GetContact(1)},
		Data: GetRequest[T]{
			QueryHost: c.GetContact(0),
			GetId:     taskId,
			Key:       key,
		},
	}
	c.SendRequestUntilConfirmation(clientTask, taskId)
	confirmation := c.GetPendingResponse(taskId)
	return confirmation.Value, confirmation.Confirmation
}

func (c *BruteChord[T]) Put(key ChordHash, value []byte) bool {
	c.logger.WriteToFileOK("Calling Put Method with key %v", key)
	// create the taskId waiting for the response and send a put request to yourself. It is the same logic for the Get now.
	taskId := generateTaskId()
	taskClient := ClientTask[T]{
		Targets: []T{c.GetContact(1)},
		Data: PutRequest[T]{
			QueryHost: c.GetContact(0),
			PutId:     taskId,
			Key:       key,
			Value:     value,
		},
	}
	c.SendRequestUntilConfirmation(taskClient, taskId)
	confirmation := c.GetPendingResponse(taskId)
	return confirmation.Confirmation
}

func (c *BruteChord[T]) sendCheckPredecessor() {
	c.logger.WriteToFileOK("Calling sendCheckPredecessor Method")
	c.logger.WriteToFileOK("Sending AreYouMyPredecessor Notification to Everyone")
	c.ClientChordCommunication.sendRequestEveryone(AreYouMyPredecessor[T]{
		Contact:     c.GetContact(0),
		MySuccessor: c.GetContact(1),
	})
}

func (c *BruteChord[T]) sendCheckAlive() {
	c.logger.WriteToFileOK("Calling sendCheckAlive Method")
	c.logger.WriteToFileOK("Sending AreYouAliveNotification to %v", c.GetContact(1))
	c.ClientChordCommunication.sendRequest(ClientTask[T]{
		Targets: []T{c.GetContact(1), c.GetContact(2)},
		Data:    AreYouAliveNotification[T]{Contact: c.GetContact(0)},
	})
}

func (c *BruteChord[T]) killDead() {
	c.logger.WriteToFileOK("Calling killDead Method at time %v", time.Now())
	successor := c.GetContact(1)
	if !c.Monitor.CheckAlive(successor, 3*WaitingTime) {
		c.logger.WriteToFileOK("Successor %v is Dead", successor)
		c.DeadContacts = append(c.DeadContacts, successor)
		c.Monitor.DeleteContact(successor)
		c.SetContact(c.DefaultSuccessor(), 1)
	} else {
		c.logger.WriteToFileOK("Successor %v is Alive", successor)
	}
	successorSuccessor := c.GetContact(2)
	if !c.Monitor.CheckAlive(successorSuccessor, 3*WaitingTime) {
		c.logger.WriteToFileOK("My Successor Successor %v looks Dead to me", successorSuccessor)
		c.DeadContacts = append(c.DeadContacts, successorSuccessor)
		c.Monitor.DeleteContact(successorSuccessor)
		c.SetContact(c.GetContact(1), 2)
	}
}

func (c *BruteChord[T]) StopWorking() {
	c.NotificationChannelServerNode <- nil
}

func sortKeys(data Store) []ChordHash {
	keys := make([]ChordHash, 0)
	for key := range data {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

func (c *BruteChord[T]) State() string {
	state := "Node: " + strconv.Itoa(int(c.GetId())) + "\n"
	state += "Successor: " + strconv.Itoa(int(c.GetContact(1).getNodeId())) + "\n"
	state += "Successor Data Replicas Are: " + "\n"

	successorData := c.GetAllData(1)
	successorSuccessorData := c.GetAllData(2)
	ownData := c.GetAllData(0)
	for _, key := range sortKeys(successorData) {
		state += strconv.Itoa(int(key)) + " -> " + fmt.Sprintf("%v", successorData[ChordHash(key)]) + "\n"
	}
	state += "SuccessorSuccessor: " + strconv.Itoa(int(c.GetContact(2).getNodeId())) + "\n"
	state += "SuccessorSuccessor Data Replica:" + "\n"
	for _, key := range sortKeys(successorSuccessorData) {
		state += strconv.Itoa(int(key)) + " -> " + fmt.Sprintf("%v", successorSuccessorData[ChordHash(key)]) + "\n"
	}
	state += "Predecessor: " + strconv.Itoa(int(c.GetContact(-1).getNodeId())) + "\n"
	state += "Data stored:\n"
	for _, key := range sortKeys(ownData) {
		state += strconv.Itoa(int(key)) + " -> " + fmt.Sprintf("%v", ownData[ChordHash(key)]) + "\n"
	}
	return state
}
