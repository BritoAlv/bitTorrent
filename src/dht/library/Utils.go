package library

import (
	"fmt"
	"math/rand"
)

var usedId = map[ChordHash]bool{}

func ConvertStr(value ChordHash) string {
	var result string
	for i := 0; i < NumberBits; i++ {
		result += string(value[i] + '0')
	}
	return result
}

func EqualBytesArray(A []byte, B []byte) bool {
	if len(A) != len(B) {
		return false
	}
	for i := 0; i < len(A); i++ {
		if A[i] != B[i] {
			return false
		}
	}
	return true
}

// Between : starting from L + 1 in a clockwise order, I can reach M before R + 1.
func Between(L ChordHash, M ChordHash, R ChordHash) bool {
	l := BinaryArrayToInt(L)
	m := BinaryArrayToInt(M)
	r := BinaryArrayToInt(R)
	fmt.Println(l, m, r)
	l = (l + 1) % (1 << NumberBits)
	for {
		if l == m {
			return true
		}
		if l == r {
			break
		}
		l = (l + 1) % (1 << NumberBits)
	}
	return false
}

func GenerateRandomBinaryId() ChordHash {
	var result ChordHash
	for {
		for i := 0; i < NumberBits; i++ {
			number := rand.Float32()
			if number >= 0.5 {
				result[i] = 1
			} else {
				result[i] = 0
			}
		}
		if _, exist := usedId[result]; !exist {
			usedId[result] = true
			break
		}
	}
	return result
}

func IntToBinaryArray(number int) ChordHash {
	var result ChordHash
	if number >= (1 << NumberBits) {
		return result
	}
	for i := 0; i < NumberBits; i++ {
		result[i] = uint8(number % 2)
		number = number / 2
	}
	return result
}

func BinaryArrayToInt(array ChordHash) int {
	result := 0
	for i := 0; i < NumberBits; i++ {
		result = result*2 + int(array[NumberBits-i-1])
	}
	return result
}
