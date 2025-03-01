package MonitorHand

import (
	"bittorrent/common"
	"bittorrent/dht/library/BruteChord"
	"bittorrent/dht/library/BruteChord/Core"
	"time"
)

type MonitorHand[T Core.Contact] struct {
	logger        common.Logger
	lastDateKnown BruteChord.SafeMap[Core.ChordHash, time.Time]
}

func NewMonitorHand[T Core.Contact](name string) *MonitorHand[T] {
	return &MonitorHand[T]{
		lastDateKnown: BruteChord.SafeMap[Core.ChordHash, time.Time]{},
		logger:        *common.NewLogger(name + ".txt"),
	}
}

func (m *MonitorHand[T]) AddContact(contact T, date time.Time) {
	m.logger.WriteToFileOK("Adding contact %v to MonitorHand at time %v", contact.GetNodeId(), date)
	m.lastDateKnown.Set(contact.GetNodeId(), date)
}

func (m *MonitorHand[T]) UpdateContactDate(contact T, date time.Time) {
	m.logger.WriteToFileOK("Updating contact %v to MonitorHand at time %v", contact.GetNodeId(), date)
	_, exist := m.lastDateKnown.Get(contact.GetNodeId())
	if exist {
		m.lastDateKnown.Set(contact.GetNodeId(), date)
	}
}

func (m *MonitorHand[T]) DeleteContact(contact T) {
	m.logger.WriteToFileOK("Deleting contact %v from MonitorHand", contact.GetNodeId())
	m.lastDateKnown.Delete(contact.GetNodeId())
}

func (m *MonitorHand[T]) CheckAlive(contact T, seconds int) bool {
	date, exist := m.lastDateKnown.Get(contact.GetNodeId())
	m.logger.WriteToFileOK("Checking if contact %v is alive", contact.GetNodeId())
	if exist {
		m.logger.WriteToFileOK("Contact %v was last seen at %v", contact.GetNodeId(), date)
	} else {
		m.logger.WriteToFileOK("Contact %v was never seen", contact.GetNodeId())
	}
	return exist && date.Add(time.Duration(seconds)*time.Second).After(time.Now())
}
