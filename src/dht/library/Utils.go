package library

import "math/rand"

func GenerateRandomBinaryId() [NumberBits]uint8 {
	var result [NumberBits]uint8
	for i := 0; i < NumberBits; i++ {
		number := rand.Float32()
		if number >= 0.5 {
			result[i] = 1
		} else {
			result[i] = 0
		}
	}
	return result
}

func IntToBinaryArray(number int) [NumberBits]uint8 {
	var result [NumberBits]uint8
	if number >= (1 << NumberBits) {
		return result
	}
	for i := 0; i < NumberBits; i++ {
		result[i] = uint8(number % 2)
		number = number / 2
	}
	return result
}
