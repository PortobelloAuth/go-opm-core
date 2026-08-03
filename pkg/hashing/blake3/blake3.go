package blake3

import (
	"fmt"

	"lukechampine.com/blake3"
)

type Blake3 struct {
	Hash *blake3.Hasher
}

// TODO: consider whether we want to implement a blake3 key derivation function version.
// NOTE: Blake3 is probably too fast and parallelizable for use in routing keys (unless
// truncation of the routing key somehow creates a sufficient number of collisions).

func NewBlake3(sz int, k []byte) (*Blake3, error) {
	// make sure that sz is a resonable size
	if sz < 1 {
		return nil, fmt.Errorf("%d is not a supported Blake3 hash size", sz)
	}

	// make sure k is nil or 32 bytes long
	if k != nil && len(k) != 32 {
		return nil, fmt.Errorf("if set, key must be 32 bytes long: %d", len(k))
	}

	return &Blake3{
		Hash: blake3.New(sz, k),
	}, nil
}

func (b *Blake3) HashFn(source []byte, salt []byte) ([]byte, error) {
	// check that salt length is 12+
	if len(salt) < 12 {
		return nil, fmt.Errorf("salt length must be at least 12: %d", len(salt))
	}

	half := len(salt) / 2
	input := append(salt[0:half], source...)
	input = append(input, salt[half:]...)

	b.Hash.Reset()
	b.Hash.Write(input)

	return b.Hash.Sum(nil), nil
}

func (b *Blake3) HashFunc(source string, salt []byte) ([]byte, error) {
	return b.HashFn([]byte(source), salt)
}
