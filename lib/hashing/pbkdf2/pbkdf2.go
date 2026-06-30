package pbkdf2

import (
	"crypto/pbkdf2"
	"crypto/sha256"
	"fmt"
	"hash"
)

// set up a constant mapping of strings to supported hash functions
// (TODO: see https://pkg.go.dev/crypto for some additional functions that we should consider)
var supportedHashes = map[string]func() hash.Hash{
	"SHA2_256": sha256.New,
}

// set up a mapping of predefined PBKDF2-{hash name}-{iterations}-{bit length} strings
var opmHashTypeDefs = map[string]PBKDF2{
	"PBKDF2-SHA2_256-100000-256": {
		Hash:       sha256.New,
		Iterations: 100000,
		KeyLength:  32, // 32 bytes is 256 bits
	},
}

type PBKDF2 struct {
	Hash       func() hash.Hash
	Iterations int
	KeyLength  int
}

func NewPBKDF2HashType(hashtype string) (*PBKDF2, error) {
	// find settings for a known hash type
	h, ok := opmHashTypeDefs[hashtype]
	if !ok {
		return nil, fmt.Errorf("%s is not a supported hash type", hashtype)
	}

	return &PBKDF2{
		Hash:       h.Hash,
		Iterations: h.Iterations,
		KeyLength:  h.KeyLength,
	}, nil
}

func NewPBKDF2(hstr string, i int, kl int) (*PBKDF2, error) {
	// make sure that Hash is a supported hash
	h, ok := supportedHashes[hstr]
	if !ok {
		return nil, fmt.Errorf("%s is not a supported hash", hstr)
	}

	// make sure that i is large enough
	if i < 10000 {
		return nil, fmt.Errorf("%d is too few iterations", i)
	}

	// test the key length against the hash
	_, err := pbkdf2.Key(h, "test", []byte("key length test salt"), 1, kl)
	if err != nil {
		return nil, err
	}

	return &PBKDF2{
		Hash:       h,
		Iterations: i,
		KeyLength:  kl,
	}, nil
}

func (p *PBKDF2) HashFn(source []byte, salt []byte) ([]byte, error) {
	// NOTE: this conversion from bytes to string may be vulnerable to cryptographic timing attacks
	// It also may not be exactly what we want?
	return p.HashFunc(string(source), salt)
}

func (p *PBKDF2) HashFunc(source string, salt []byte) ([]byte, error) {
	// check that salt length is 12+
	if len(salt) < 12 {
		return nil, fmt.Errorf("salt length must be at least 12: %d", len(salt))
	}

	return pbkdf2.Key(p.Hash, source, salt, p.Iterations, p.KeyLength)
}
