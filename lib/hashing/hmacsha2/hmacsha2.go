package hmacsha2

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
)

type SHA2 struct {
	Hash func() hash.Hash
}

func NewHMACSHA2(name string) (*SHA2, error) {
	// make sure that Hash is a supported hash
	switch name {
	case "224":
		return &SHA2{
			Hash: sha256.New224,
		}, nil
	case "256":
		return &SHA2{
			Hash: sha256.New,
		}, nil
	case "384":
		return &SHA2{
			Hash: sha512.New384,
		}, nil
	case "512":
		return &SHA2{
			Hash: sha512.New,
		}, nil
	case "512_224":
		return &SHA2{
			Hash: sha512.New512_224,
		}, nil
	case "512_256":
		return &SHA2{
			Hash: sha512.New512_256,
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

	h := hmac.New(s.Hash, salt)
	h.Reset()
	h.Write(source)

	return h.Sum(nil), nil
}

func (s *SHA2) HashFunc(source string, salt []byte) ([]byte, error) {
	return s.HashFn([]byte(source), salt)
}

// TODO: consider using use hmac's constant time comparison to mitigate timing
// attacks - though it is unclear how much that matters in a system that delivers
// the hash over a network (and therefore doesn't control the client.)
