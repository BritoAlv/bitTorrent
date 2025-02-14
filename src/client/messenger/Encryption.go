package messenger

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"hash"
	"io"
)

func encrypt(message []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	encryptedMessage, err := encryptOAEP(sha256.New(), rand.Reader, publicKey, message, nil)
	if err != nil {
		return nil, errors.New("error while encrypting message")
	}
	return encryptedMessage, nil
}

func decrypt(message []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	decryptedMessage, err := decryptOAEP(sha256.New(), rand.Reader, privateKey, message, nil)
	if err != nil {
		return nil, errors.New("error while decrypting message")
	}
	return decryptedMessage, nil
}

func encryptOAEP(hash hash.Hash, random io.Reader, public *rsa.PublicKey, msg []byte, label []byte) ([]byte, error) {
	msgLen := len(msg)
	step := public.Size() - 2*hash.Size() - 2
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(hash, random, public, msg[start:finish], label)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

func decryptOAEP(hash hash.Hash, random io.Reader, private *rsa.PrivateKey, msg []byte, label []byte) ([]byte, error) {
	msgLen := len(msg)
	step := private.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(hash, random, private, msg[start:finish], label)
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}
