package Core

type Notification[contact Contact] interface {
	HandleNotification(*BruteChord[contact])
}
