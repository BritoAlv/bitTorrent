package library

type Notification[contact Contact] interface {
	HandleNotification(*BruteChord[contact])
}
