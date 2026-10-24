package sha2

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
)

type SHA2 struct {
	Hash hash.Hash
}

// set up a mapping of predefined SHA2-{name} strings
var opmHashTypeDefs = map[string]string{
	"SHA2-224":     "224",
	"SHA2-256":     "256",
	"SHA2-384":     "384",
	"SHA2-512":     "512",
	"SHA2-512_224": "512_224",
	"SHA2-512_256": "512_256",
}

func NewSHA2HashType(hashtype string) (*SHA2, error) {
	// find name for a known hash type
	name, ok := opmHashTypeDefs[hashtype]
	if !ok {
		return nil, fmt.Errorf("%s is not a supported hash type", hashtype)
	}

	return NewSHA2(name)
}

func NewSHA2(name string) (*SHA2, error) {
	// make sure that Hash is a supported hash
	switch name {
	case "224":
		return &SHA2{
			Hash: sha256.New224(),
		}, nil
	case "256":
		return &SHA2{
			Hash: sha256.New(),
		}, nil
	case "384":
		return &SHA2{
			Hash: sha512.New384(),
		}, nil
	case "512":
		return &SHA2{
			Hash: sha512.New(),
		}, nil
	case "512_224":
		return &SHA2{
			Hash: sha512.New512_224(),
		}, nil
	case "512_256":
		return &SHA2{
			Hash: sha512.New512_256(),
		}, nil
	default:
		return nil, fmt.Errorf("%s is not a supported SHA2 hash type", name)
	}
}

func (s *SHA2) HashFn(source []byte, salt []byte) ([]byte, error) {
	// check that salt length is 12+
	if len(salt) < 12 {
		return nil, fmt.Errorf("salt length must be at least 12: %d", len(salt))
	}

	half := len(salt) / 2
	input := make([]byte, 0, len(salt)+len(source))
	input = append(input, salt[:half]...)
	input = append(input, source...)
	input = append(input, salt[half:]...)

	h := s.Hash
	h.Reset()
	h.Write(input)

	return h.Sum(nil), nil
}

func (s *SHA2) HashFunc(source string, salt []byte) ([]byte, error) {
	return s.HashFn([]byte(source), salt)
}
