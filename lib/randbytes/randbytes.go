package randbytes

import (
	"crypto/rand"
	"fmt"
)

// RandomBytes creates a cryptographically secure random byte slice.
// It is recommended to use at least 16 bytes for a secure salt.
func RandomBytes(num uint32) ([]byte, error) {
	bytes := make([]byte, num)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// GenerateSalt creates a cryptographically secure random byte slice.
// It is recommended to use at least 16 bytes for a secure salt. Therefore
// this function requires a salt length of at least 16. Otherwise it is
// just a wrapper around RandomBytes.
func GenerateSalt(length uint32) ([]byte, error) {
	if length < 16 {
		return nil, fmt.Errorf("Cryptographic salts should generally contain at least 16 bytes")
	}
	return RandomBytes(length)
}
