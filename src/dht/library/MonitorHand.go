package library

import "time"

type MonitorHand[T Contact] struct {
	lastDateKnown map[[NumberBits]uint8]time.Time
}

func NewMonitorHand[T Contact]() *MonitorHand[T] {
	return &MonitorHand[T]{
		lastDateKnown: make(map[[8]uint8]time.Time),
	}
}

func (m *MonitorHand[T]) AddContact(contact T, date time.Time) {
	m.lastDateKnown[contact.getNodeId()] = date
}

func (m *MonitorHand[T]) UpdateContactDate(contact T, date time.Time) {
	_, exist := m.lastDateKnown[contact.getNodeId()]
	if exist {
		m.lastDateKnown[contact.getNodeId()] = date
	}
}

func (m *MonitorHand[T]) DeleteContact(contact T) {
	delete(m.lastDateKnown, contact.getNodeId())
}

func (m *MonitorHand[T]) CheckAlive(contact T, second int) bool {
	date, exist := m.lastDateKnown[contact.getNodeId()]
	return exist && date.Add(time.Duration(second)).After(time.Now())
}
