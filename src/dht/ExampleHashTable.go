package main

import "bittorrent/dht/library"

type CustomHashTable[T any, U library.IHashable] struct {
	dictionary map[[library.NumberBits]uint8]T
}

func (d *CustomHashTable[T, U]) Clear() {
	d.dictionary = map[[library.NumberBits]uint8]T{}
}

func NewCustomHashTable[T any, U library.IHashable]() *CustomHashTable[T, U] {
	return &CustomHashTable[T, U]{dictionary: map[[library.NumberBits]uint8]T{}}
}

func (d *CustomHashTable[T, U]) Put(key U, value T) {
	d.dictionary[key.Hash()] = value
}

func (d *CustomHashTable[T, U]) Get(key U) (T, bool) {
	result, exist := d.dictionary[key.Hash()]
	return result, exist
}
