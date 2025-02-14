package library

type ExampleHashTable[T any, U IHashable] struct {
	dictionary map[[NumberBits]uint8]T
}

func (d *ExampleHashTable[T, U]) Clear() {
	d.dictionary = map[[NumberBits]uint8]T{}
}

func NewExampleHashTable[T any, U IHashable]() *ExampleHashTable[T, U] {
	return &ExampleHashTable[T, U]{dictionary: map[[NumberBits]uint8]T{}}
}

func (d *ExampleHashTable[T, U]) Put(key U, value T) {
	d.dictionary[key.Hash()] = value
}

func (d *ExampleHashTable[T, U]) Get(key U) (T, bool) {
	result, exist := d.dictionary[key.Hash()]
	return result, exist
}
