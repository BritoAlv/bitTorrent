package messenger

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
)

func encrypt(message []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	encryptedMessage, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, message, nil)
	if err != nil {
		return nil, errors.New("Error while encrypting message")
	}
	return encryptedMessage, nil
}

func decrypt(message []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	decryptedMessage, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, message, nil)
	if err != nil {
		return nil, errors.New("Error while decrypting message")
	}
	return decryptedMessage, nil
}
