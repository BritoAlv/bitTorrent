package library

import (
	"bittorrent/common"
	"fmt"
	"time"
)

type MonitorHand[T Contact] struct {
	logger        common.Logger
	lastDateKnown map[[NumberBits]uint8]time.Time
}

func NewMonitorHand[T Contact]() *MonitorHand[T] {
	return &MonitorHand[T]{
		lastDateKnown: make(map[[8]uint8]time.Time),
		logger:        *common.NewLogger("MonitorHand.txt"),
	}
}

func (m *MonitorHand[T]) AddContact(contact T, date time.Time) {
	m.logger.WriteToFileOK(fmt.Sprintf("Adding contact %v to MonitorHand at time %v", contact.getNodeId(), date))
	m.lastDateKnown[contact.getNodeId()] = date
}

func (m *MonitorHand[T]) UpdateContactDate(contact T, date time.Time) {
	m.logger.WriteToFileOK(fmt.Sprintf("Updating contact %v to MonitorHand at time %v", contact.getNodeId(), date))
	_, exist := m.lastDateKnown[contact.getNodeId()]
	if exist {
		m.lastDateKnown[contact.getNodeId()] = date
	}
}

func (m *MonitorHand[T]) DeleteContact(contact T) {
	m.logger.WriteToFileOK(fmt.Sprintf("Deleting contact %v from MonitorHand", contact.getNodeId()))
	delete(m.lastDateKnown, contact.getNodeId())
}

func (m *MonitorHand[T]) CheckAlive(contact T, seconds int) bool {
	date, exist := m.lastDateKnown[contact.getNodeId()]
	m.logger.WriteToFileOK(fmt.Sprintf("Checking if contact %v is alive", contact.getNodeId()))
	if exist {
		m.logger.WriteToFileOK(fmt.Sprintf("Contact %v was last seen at %v", contact.getNodeId(), date))
	} else {
		m.logger.WriteToFileOK(fmt.Sprintf("Contact %v was never seen", contact.getNodeId()))
	}
	return exist && date.Add(time.Duration(seconds)*time.Second).After(time.Now())
}
