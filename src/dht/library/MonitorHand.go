package library

import (
	"bittorrent/common"
	"sync"
	"time"
)

type MonitorHand[T Contact] struct {
	logger        common.Logger
	lastDateKnown map[ChordHash]time.Time
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
	m.logger.WriteToFileOK("Adding contact %v to MonitorHand at time %v", contact.getNodeId(), date)
	m.lastDateKnown[contact.getNodeId()] = date
	m.lock.Unlock()
}

func (m *MonitorHand[T]) UpdateContactDate(contact T, date time.Time) {
	m.lock.Lock()
	m.logger.WriteToFileOK("Updating contact %v to MonitorHand at time %v", contact.getNodeId(), date)
	_, exist := m.lastDateKnown[contact.getNodeId()]
	if exist {
		m.lastDateKnown[contact.getNodeId()] = date
	}
	m.lock.Unlock()
}

func (m *MonitorHand[T]) DeleteContact(contact T) {
	m.lock.Lock()
	m.logger.WriteToFileOK("Deleting contact %v from MonitorHand", contact.getNodeId())
	delete(m.lastDateKnown, contact.getNodeId())
	m.lock.Unlock()
}

func (m *MonitorHand[T]) CheckAlive(contact T, seconds int) bool {
	m.lock.Lock()
	date, exist := m.lastDateKnown[contact.getNodeId()]
	m.lock.Unlock()
	m.logger.WriteToFileOK("Checking if contact %v is alive", contact.getNodeId())
	if exist {
		m.logger.WriteToFileOK("Contact %v was last seen at %v", contact.getNodeId(), date)
	} else {
		m.logger.WriteToFileOK("Contact %v was never seen", contact.getNodeId())
	}
	return exist && date.Add(time.Duration(seconds)*time.Second).After(time.Now())
}
