package Core

import (
	"time"
)

type Monitor[T Contact] interface {
	// AddContact :Add a contact to the monitor using arrivalDate as reference.
	AddContact(contact T, arrivalDate time.Time)
	// UpdateContactDate :Update the date of the contact. Users of the monitor are responsible for providing the date.
	UpdateContactDate(contact T, date time.Time)
	// DeleteContact :Delete a contact from the monitor.
	DeleteContact(contact T)
	// CheckAlive :A contact is alive if it has been updated in the last seconds.
	CheckAlive(contact T, seconds int) bool
}
