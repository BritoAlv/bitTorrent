package library

import (
	"fmt"
	rand2 "math/rand/v2"
)

var usedId = map[ChordHash]bool{}

func GenerateRandomBinaryId() ChordHash {
	var result ChordHash
	for {
		result = rand2.Int64() % (1 << NumberBits)
		if _, exist := usedId[result]; !exist {
			usedId[result] = true
			break
		}
	}
	return result
}

// Between : starting from L + 1 in a clockwise order, I can reach M before R + 1.
func Between(L ChordHash, M ChordHash, R ChordHash) bool {
	fmt.Println(L, M, R)
	L = (L + 1) % (1 << NumberBits)
	for {
		if L == M {
			return true
		}
		if L == R {
			break
		}
		L = (L + 1) % (1 << NumberBits)
	}
	return false
}
