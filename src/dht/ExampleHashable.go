package main

import "bittorrent/dht/library"

type ExampleHashable struct {
	data [library.NumberBits]uint8
}

func (e ExampleHashable) Hash() [library.NumberBits]uint8 {
	return e.data
}

func NewExampleHashable(number int) *ExampleHashable {
	return &ExampleHashable{data: library.IntToBinaryArray(number)}
}
