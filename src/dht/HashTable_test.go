package main

import (
	"bittorrent/dht/library"
	"math/rand"
	"testing"
)

var tableOk = library.NewExampleHashTable[int, library.ExampleHashable]()
var tableTest = library.NewExampleHashTable[int, library.ExampleHashable]()

func ClearTable() {
	tableTest.Clear()
	tableOk.Clear()
}

func TestHashTablePut(t *testing.T) {
	ClearTable()
	for i := 0; i < 100; i++ {
		key := *library.NewExampleHashable(rand.Int() % (1 << library.NumberBits))
		value := rand.Int()
		tableOk.Put(key, value)
		tableTest.Put(key, value)
	}
	for i := 0; i < 100; i++ {
		keyToFind := *library.NewExampleHashable(i)
		valueOk, existOk := tableOk.Get(keyToFind)
		valueTest, existTest := tableTest.Get(keyToFind)
		if existOk != existTest {
			t.Errorf("When doing the get for key %v, the exist value is different: ok table says %v, different from %v", keyToFind, existOk, existTest)
		}
		if existOk && existTest {
			if valueOk != valueTest {
				t.Errorf("When doing the get for key %v, the value is different: ok table says %v, different from %v", keyToFind, valueOk, valueTest)
			}
		}
	}
}

func TestHashTableFalseGet(t *testing.T) {
	ClearTable()
	key := *library.NewExampleHashable(2)
	_, exist := tableTest.Get(key)
	if exist {
		t.Errorf("Expected not to find the value")
	}
}
