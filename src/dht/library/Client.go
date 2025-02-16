package library

type Client[T Contact] interface {
	sendRequest(task ClientTask[T])
	sendRequestEveryone(data Notification[T])
}
type ClientTask[T Contact] struct {
	Targets []T
	Data    Notification[T]
}
