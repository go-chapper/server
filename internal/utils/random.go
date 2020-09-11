package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// RandomCryptoString returns a cryptographically secure random string
func RandomCryptoString(n int) (string, error) {
	s := make([]byte, n)
	_, err := rand.Read(s)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(s), nil
}

func RandomByteSlice(n uint32) ([]byte, error) {
	s := make([]byte, n)
	_, err := rand.Read(s)
	if err != nil {
		return nil, err
	}

	return s, nil
}
