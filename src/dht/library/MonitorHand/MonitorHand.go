package MonitorHand

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord/Core"
	"sync"
	"time"
)

type MonitorHand[T Core.Contact] struct {
	logger        common.Logger
	lastDateKnown map[Core.ChordHash]time.Time
	lock          sync.Mutex
}

func NewMonitorHand[T Core.Contact](name string) *MonitorHand[T] {
	return &MonitorHand[T]{
		lastDateKnown: make(map[Core.ChordHash]time.Time),
		logger:        *common.NewLogger(name + ".txt"),
		lock:          sync.Mutex{},
	}
}

func (m *MonitorHand[T]) AddContact(contact T, date time.Time) {
	m.lock.Lock()
	m.logger.WriteToFileOK("Adding contact %v to MonitorHand at time %v", contact.GetNodeId(), date)
	m.lastDateKnown[contact.GetNodeId()] = date
	m.lock.Unlock()
}

func (m *MonitorHand[T]) UpdateContactDate(contact T, date time.Time) {
	m.lock.Lock()
	m.logger.WriteToFileOK("Updating contact %v to MonitorHand at time %v", contact.GetNodeId(), date)
	_, exist := m.lastDateKnown[contact.GetNodeId()]
	if exist {
		m.lastDateKnown[contact.GetNodeId()] = date
	}
	m.lock.Unlock()
}

func (m *MonitorHand[T]) DeleteContact(contact T) {
	m.lock.Lock()
	m.logger.WriteToFileOK("Deleting contact %v from MonitorHand", contact.GetNodeId())
	delete(m.lastDateKnown, contact.GetNodeId())
	m.lock.Unlock()
}

func (m *MonitorHand[T]) CheckAlive(contact T, seconds int) bool {
	m.lock.Lock()
	date, exist := m.lastDateKnown[contact.GetNodeId()]
	m.lock.Unlock()
	m.logger.WriteToFileOK("Checking if contact %v is alive", contact.GetNodeId())
	if exist {
		m.logger.WriteToFileOK("Contact %v was last seen at %v", contact.GetNodeId(), date)
	} else {
		m.logger.WriteToFileOK("Contact %v was never seen", contact.GetNodeId())
	}
	return exist && date.Add(time.Duration(seconds)*time.Second).After(time.Now())
}
