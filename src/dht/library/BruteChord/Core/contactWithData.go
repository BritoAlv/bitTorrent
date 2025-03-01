package Core

type contactWithData[T Contact] struct {
	Contact T
	Data    SafeStore
}
