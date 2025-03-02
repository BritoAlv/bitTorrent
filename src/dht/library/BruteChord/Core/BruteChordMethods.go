package Core

import (
	"fmt"
	"strconv"
	"time"
)

func (c *BruteChord[T]) responsible(key ChordHash) bool {
	return between(c.GetId(), key, c.GetContact(1).GetNodeId())
}

func (c *BruteChord[T]) sendRequestUntilConfirmation(clientTask ClientTask[T], taskId int64) bool {
	c.logger.WriteToFileOK("Calling sendRequestUntilConfirmation Method with request %v and taskId %v", clientTask, taskId)
	c.setPendingResponse(taskId, confirmations{Confirmation: false, Value: nil})
	for i := 0; i < Attempts; i++ {
		confirmation := c.getPendingResponse(taskId)
		if confirmation.Confirmation {
			c.logger.WriteToFileOK("Received Confirmation for taskId %v", taskId)
			return true
		}
		c.clientChordCommunication.SendRequest(clientTask)
		time.Sleep(WaitingTime * time.Second)
	}
	c.logger.WriteToFileOK("Didn't receive Confirmation for taskId %v", taskId)
	return false
}

func (c *BruteChord[T]) replicateData(successorIndex int, target T) {
	c.logger.WriteToFileOK("Calling replicateData Method with successorIndex %v designated to %v", successorIndex, target)
	taskId := generateTaskId()
	clientTask := ClientTask[T]{
		Targets: []T{target},
		Data: receiveDataReplicate[T]{
			SuccessorIndex: successorIndex,
			TaskId:         taskId,
			Data:           c.getAllData(0),
			DataOwner:      c.GetContact(0),
		},
	}
	c.sendRequestUntilConfirmation(clientTask, taskId)
}

func (c *BruteChord[T]) Get(key ChordHash) ([]byte, bool) {
	c.logger.WriteToFileOK("Calling Get Method on key %v", key)
	taskId := generateTaskId()
	clientTask := ClientTask[T]{
		Targets: []T{c.GetContact(0)},
		Data: getRequest[T]{
			QueryHost: c.GetContact(0),
			GetId:     taskId,
			Key:       key,
		},
	}
	c.sendRequestUntilConfirmation(clientTask, taskId)
	confirmation := c.getPendingResponse(taskId)
	if confirmation.Value == nil {
		return nil, false
	} else {
		return confirmation.Value, true
	}
}

func (c *BruteChord[T]) Put(key ChordHash, value []byte) bool {
	c.logger.WriteToFileOK("Calling Put Method with key %v", key)
	// create the taskId waiting for the response and send a put request to yourself. It is the same logic for the Get now.
	taskId := generateTaskId()
	fmt.Printf("TaskId for the Put is %v\n", taskId)
	taskClient := ClientTask[T]{
		Targets: []T{c.GetContact(0)},
		Data: putRequest[T]{
			QueryHost: c.GetContact(0),
			PutId:     taskId,
			Key:       key,
			Value:     value,
		},
	}
	result := c.sendRequestUntilConfirmation(taskClient, taskId)
	fmt.Printf("Result %v\n", result)
	confirmation := c.getPendingResponse(taskId)
	return confirmation.Confirmation
}

func (c *BruteChord[T]) sendCheckPredecessor() {
	c.logger.WriteToFileOK("Calling sendCheckPredecessor Method")
	c.logger.WriteToFileOK("Sending areYouMyPredecessor Notification to Everyone")
	c.clientChordCommunication.SendRequestEveryone(areYouMyPredecessor[T]{
		Contact:     c.GetContact(0),
		MySuccessor: c.GetContact(1),
	})
}

func (c *BruteChord[T]) sendCheckAlive() {
	c.logger.WriteToFileOK("Calling sendCheckAlive Method")
	c.logger.WriteToFileOK("Sending areYouAliveNotification to %v", c.GetContact(1))
	c.clientChordCommunication.SendRequest(ClientTask[T]{
		Targets: []T{c.GetContact(1), c.GetContact(2), c.GetContact(-1)},
		Data:    areYouAliveNotification[T]{Contact: c.GetContact(0)},
	})
}

func (c *BruteChord[T]) killDead() {
	c.logger.WriteToFileOK("Calling killDead Method at time %v", time.Now())
	successor := c.GetContact(1)
	if !c.monitor.CheckAlive(successor, 3*WaitingTime) {
		c.logger.WriteToFileOK("Successor %v is Dead", successor)
		c.monitor.DeleteContact(successor)
		c.setContact(c.defaultSuccessor(), 1)
	} else {
		c.logger.WriteToFileOK("Successor %v is Alive", successor)
	}
	successorSuccessor := c.GetContact(2)
	if !c.monitor.CheckAlive(successorSuccessor, 3*WaitingTime) {
		c.logger.WriteToFileOK("My Successor Successor %v looks Dead to me", successorSuccessor)
		c.monitor.DeleteContact(successorSuccessor)
		c.setContact(c.GetContact(1), 2)
	}
	predecessor := c.GetContact(-1)
	if !c.monitor.CheckAlive(predecessor, 3*WaitingTime) {
		c.logger.WriteToFileOK("My Predecessor %v looks Dead to me", predecessor)
		c.monitor.DeleteContact(predecessor)
		c.setContact(c.GetContact(0), -1)
	}
}

func (c *BruteChord[T]) GetState() string {
	state := "Node: " + strconv.Itoa(int(c.GetId())) + "\n"
	state += "Successor: " + strconv.Itoa(int(c.GetContact(1).GetNodeId())) + "\n"
	state += "Successor Data Replicas Are: " + "\n"

	successorData := c.getAllData(1)
	successorSuccessorData := c.getAllData(2)
	ownData := c.getAllData(0)
	for _, key := range sortKeys(successorData) {
		state += strconv.Itoa(int(key)) + " -> " + fmt.Sprintf("%v", decodePeerList(successorData[key])) + "\n"
	}
	state += "SuccessorSuccessor: " + strconv.Itoa(int(c.GetContact(2).GetNodeId())) + "\n"
	state += "SuccessorSuccessor Data Replica:" + "\n"
	for _, key := range sortKeys(successorSuccessorData) {
		state += strconv.Itoa(int(key)) + " -> " + fmt.Sprintf("%v", decodePeerList(successorSuccessorData[key])) + "\n"
	}
	state += "Predecessor: " + strconv.Itoa(int(c.GetContact(-1).GetNodeId())) + "\n"
	state += "Data stored:\n"
	for _, key := range sortKeys(ownData) {
		state += strconv.Itoa(int(key)) + " -> " + fmt.Sprintf("%v", decodePeerList(ownData[key])) + "\n"
	}
	return state
}
