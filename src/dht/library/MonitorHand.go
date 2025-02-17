package library

import (
	"bittorrent/common"
	"fmt"
	"sync"
	"time"
)

type MonitorHand[T Contact] struct {
	logger        common.Logger
	lastDateKnown map[[NumberBits]uint8]time.Time
	lock          sync.Mutex
}

func NewMonitorHand[T Contact](name string) *MonitorHand[T] {
	return &MonitorHand[T]{
		lastDateKnown: make(map[[8]uint8]time.Time),
		logger:        *common.NewLogger(name + ".txt"),
		lock:          sync.Mutex{},
	}
}

func (m *MonitorHand[T]) AddContact(contact T, date time.Time) {
	m.lock.Lock()
	m.logger.WriteToFileOK(fmt.Sprintf("Adding contact %v to MonitorHand at time %v", contact.getNodeId(), date))
	m.lastDateKnown[contact.getNodeId()] = date
	m.lock.Unlock()
}

func (m *MonitorHand[T]) UpdateContactDate(contact T, date time.Time) {
	m.lock.Lock()
	m.logger.WriteToFileOK(fmt.Sprintf("Updating contact %v to MonitorHand at time %v", contact.getNodeId(), date))
	_, exist := m.lastDateKnown[contact.getNodeId()]
	if exist {
		m.lastDateKnown[contact.getNodeId()] = date
	}
	m.lock.Unlock()
}

func (m *MonitorHand[T]) DeleteContact(contact T) {
	m.lock.Lock()
	m.logger.WriteToFileOK(fmt.Sprintf("Deleting contact %v from MonitorHand", contact.getNodeId()))
	delete(m.lastDateKnown, contact.getNodeId())
	m.lock.Unlock()
}

func (m *MonitorHand[T]) CheckAlive(contact T, seconds int) bool {
	m.lock.Lock()
	date, exist := m.lastDateKnown[contact.getNodeId()]
	m.lock.Unlock()
	m.logger.WriteToFileOK(fmt.Sprintf("Checking if contact %v is alive", contact.getNodeId()))
	if exist {
		m.logger.WriteToFileOK(fmt.Sprintf("Contact %v was last seen at %v", contact.getNodeId(), date))
	} else {
		m.logger.WriteToFileOK(fmt.Sprintf("Contact %v was never seen", contact.getNodeId()))
	}
	return exist && date.Add(time.Duration(seconds)*time.Second).After(time.Now())
}
