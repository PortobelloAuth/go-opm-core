package randbytes

import (
	"crypto/rand"
	"fmt"
	"math"
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

func RandomBits(length uint32) ([]byte, error) {
	bytelength := uint32(math.Ceil(float64(length) / 8))

	bytes, err := RandomBytes(bytelength)
	if err != nil {
		return nil, fmt.Errorf("error generating random bits: %w", err)
	}

	bytes, err = MaskToBits(bytes, length)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func MaskToBits(source []byte, length uint32) ([]byte, error) {
	bytelength := uint32(math.Ceil(float64(length) / 8))
	if bytelength > uint32(len(source)) {
		return nil, fmt.Errorf("%d bits requires %d bytes; source has %d", length, bytelength, len(source))
	}
	bytes := source[0:bytelength]

	// zero out the excess bits
	excess := bytelength*8 - length
	if excess > 0 {
		lastbyte := bytes[len(bytes)-1]
		// Start with 1 byte of all 1's and then shift 0's on from the left
		// (because a byte is a uint8) which will result in keeping only the
		// low-order bits that are not extra.
		mask := byte(0xff) >> excess
		bytes[len(bytes)-1] = lastbyte & mask
	}

	return bytes, nil
}
