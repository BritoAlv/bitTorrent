package common

import (
	"crypto/rand"
)

func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		b[i] = '0' + (b[i] % 2)
	}
	return string(b), nil
}
