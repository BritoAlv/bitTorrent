package library

type IHashable interface {
	Hash() [NumberBits]uint8
}

type HashTable[T any, U IHashable] interface {
	Put(key U, value T)
	Get(key U) (T, bool)
	Clear()
}
