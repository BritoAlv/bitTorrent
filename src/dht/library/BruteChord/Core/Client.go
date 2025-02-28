package Core

type Client[T Contact] interface {
	SendRequest(task ClientTask[T])
	SendRequestEveryone(data Notification[T])
}
type ClientTask[T Contact] struct {
	Targets []T
	Data    Notification[T]
}
