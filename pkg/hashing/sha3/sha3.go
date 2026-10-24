package sha3

import (
	"crypto/sha3"
	"fmt"
)

var sha3Constructors = map[string]func() *sha3.SHA3{
	"224": sha3.New224,
	"256": sha3.New256,
	"384": sha3.New384,
	"512": sha3.New512,
}

type SHA3 struct {
	Hash *sha3.SHA3
}

// set up a mapping of predefined SHA3-{bit length} strings
var opmHashTypeDefs = map[string]string{
	"SHA3-224": "224",
	"SHA3-256": "256",
	"SHA3-384": "384",
	"SHA3-512": "512",
}

func NewSHA3HashType(hashtype string) (*SHA3, error) {
	// find name for a known hash type
	length, ok := opmHashTypeDefs[hashtype]
	if !ok {
		return nil, fmt.Errorf("%s is not a supported hash type", hashtype)
	}

	return NewSHA3(length)
}

func NewSHA3(length string) (*SHA3, error) {
	// make sure that Hash is a supported hash
	con, ok := sha3Constructors[length]
	if !ok {
		return nil, fmt.Errorf("%s is not a supported SHA3 hash length", length)
	}

	return &SHA3{
		Hash: con(),
	}, nil
}

func (h *SHA3) HashFn(source []byte, salt []byte) ([]byte, error) {
	// check that salt length is 12+
	if len(salt) < 12 {
		return nil, fmt.Errorf("salt length must be at least 12: %d", len(salt))
	}

	half := len(salt) / 2
	input := make([]byte, 0, len(salt)+len(source))
	input = append(input, salt[:half]...)
	input = append(input, source...)
	input = append(input, salt[half:]...)

	h.Hash.Reset()
	h.Hash.Write(input)

	return h.Hash.Sum(nil), nil
}

func (h *SHA3) HashFunc(source string, salt []byte) ([]byte, error) {
	return h.HashFn([]byte(source), salt)
}
