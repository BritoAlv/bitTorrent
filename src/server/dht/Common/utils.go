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

func BuildIPRequest(serverUrl string, name string) string {
	return serverUrl + IPRoute + "?name=" + name
}
func BuildLoginRequest(serverUrl string, name string, ip string) string {
	return serverUrl + LoginRoute + "?name=" + name + "&ip=" + ip
}
