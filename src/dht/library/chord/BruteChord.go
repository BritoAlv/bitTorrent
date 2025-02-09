package chord

import "bittorrent/dht/library"

type BruteChord[T any, U library.IHashable] struct {
	Id    [library.NumberBits]uint8
	Store library.ExampleHashTable[T, U]
}

func NewBruteChord[T any, U library.IHashable]() *BruteChord[T, U] {
	Id := library.GenerateRandomBinaryId()
	Store := *library.NewExampleHashTable[T, U]()
	return &BruteChord[T, U]{Id, Store}
}

func (b *BruteChord[T, U]) Put(key U, value T) {
	//TODO implement me
	panic("implement me")
}

func (b *BruteChord[T, U]) Get(key U) (T, bool) {
	//TODO implement me
	panic("implement me")
}

func (b *BruteChord[T, U]) Clear() {
	b.Store.Clear()
}
