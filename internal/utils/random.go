package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// RandomCryptoString returns a cryptographically secure random string with length n
func RandomCryptoString(n int) (string, error) {
	s := make([]byte, n)
	_, err := rand.Read(s)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(s), nil
}

// RandomString returns a random string with length n
func RandomString(n int) (string, error) {
	b, err := RandomByteSlice(uint32(n))
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// RandomByteSlice returns a random byte slice with length n
func RandomByteSlice(n uint32) ([]byte, error) {
	s := make([]byte, n)
	_, err := rand.Read(s)
	if err != nil {
		return nil, err
	}

	return s, nil
}
