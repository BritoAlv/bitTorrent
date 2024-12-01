package common

import (
	"errors"
	"math/rand"
)

func GenerateRandomString(size int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, size)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GetTotalPieces(length int, pieceLength int) int {
	var totalPieces int

	if length <= pieceLength {
		totalPieces = 1
	} else {
		totalPieces = length / pieceLength
		if length%pieceLength != 0 {
			totalPieces += 1
		}
	}
	return totalPieces
}

func CastTo[T interface{}](obj interface{}) (T, error) {
	casted, isExpectedType := obj.(T)

	if !isExpectedType {
		return casted, errors.New("cannot cast the object to the given type")
	}

	return casted, nil
}
