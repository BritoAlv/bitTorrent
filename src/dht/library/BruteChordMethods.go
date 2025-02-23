package library

import (
	"time"
)

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
			Data:           c.GetAllOwnData(),
			DataOwner:      c.GetContact(),
		},
	}
	go c.SendRequestUntilConfirmation(clientTask, taskId)
}

func (c *BruteChord[T]) Get(key ChordHash) []byte {
	c.logger.WriteToFileOK("Calling Get Method on key %v", key)
	taskId := generateTaskId()
	clientTask := ClientTask[T]{
		Targets: []T{c.GetSuccessor()},
		Data: GetRequest[T]{
			QueryHost: c.GetContact(),
			GetId:     taskId,
			Key:       key,
		},
	}
	c.SendRequestUntilConfirmation(clientTask, taskId)
	confirmation := c.GetPendingResponse(taskId)
	return confirmation.Value
}

func (c *BruteChord[T]) Put(key ChordHash, value []byte) bool {
	c.logger.WriteToFileOK("Calling Put Method with key %v", key)
	// create the taskId waiting for the response and send a put request to yourself. It is the same logic for the Get now.
	taskId := generateTaskId()
	taskClient := ClientTask[T]{
		Targets: []T{c.GetSuccessor()},
		Data: PutRequest[T]{
			QueryHost: c.GetContact(),
			PutId:     taskId,
			Key:       key,
			Value:     value,
		},
	}
	c.SendRequestUntilConfirmation(taskClient, taskId)
	confirmation := c.GetPendingResponse(taskId)
	return confirmation.Confirmation
}

func (c *BruteChord[T]) StabilizeStore() {
	curDate := time.Now()
	c.lock.Lock()
	c.logger.WriteToFileOK("Stabilizing Store at time %v", curDate)
	ownStore := make(map[ChordHash][]byte)
	for key, data := range c.myData {
		ownStore[key] = data
	}
	c.myData = make(map[ChordHash][]byte)
	for key, data := range ownStore {
		go c.Put(key, data)
	}
	c.lock.Unlock()
}

func (c *BruteChord[T]) sendCheckPredecessor() {
	c.logger.WriteToFileOK("Calling sendCheckPredecessor Method")
	c.logger.WriteToFileOK("Sending AreYouMyPredecessor Notification to Everyone")
	c.ClientChordCommunication.sendRequestEveryone(AreYouMyPredecessor[T]{
		Contact:     c.GetContact(),
		MySuccessor: c.GetSuccessor(),
	})
}

func (c *BruteChord[T]) sendCheckAlive() {
	c.logger.WriteToFileOK("Calling sendCheckAlive Method")
	c.logger.WriteToFileOK("Sending AreYouAliveNotification to %v", c.GetSuccessor())
	c.ClientChordCommunication.sendRequest(ClientTask[T]{
		Targets: []T{c.GetSuccessor(), c.successorSuccessor.Contact},
		Data:    AreYouAliveNotification[T]{Contact: c.GetContact()},
	})
}

func (c *BruteChord[T]) killDead() {
	c.logger.WriteToFileOK("Calling killDead Method at time %v", time.Now())
	successor := c.GetSuccessor()
	if !c.Monitor.CheckAlive(successor, 3*WaitingTime) {
		c.logger.WriteToFileOK("Successor %v is Dead", successor)
		c.DeadContacts = append(c.DeadContacts, successor)
		c.Monitor.DeleteContact(successor)
		c.SetSuccessor(c.DefaultSuccessor())
	} else {
		c.logger.WriteToFileOK("Successor %v is Alive", successor)
	}
	successorSuccessor := c.GetSuccessorSuccessor()
	if !c.Monitor.CheckAlive(successorSuccessor, 3*WaitingTime) {
		c.logger.WriteToFileOK("My Successor Successor %v looks Dead to me", successorSuccessor)
		c.DeadContacts = append(c.DeadContacts, successorSuccessor)
		c.Monitor.DeleteContact(successorSuccessor)
		c.SetSuccessorSuccessor(c.GetSuccessor())
	}
}

func (c *BruteChord[T]) StopWorking() {
	c.NotificationChannelServerNode <- nil
}
