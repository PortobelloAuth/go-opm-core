package argon2id

import (
	"fmt"

	"golang.org/x/crypto/argon2"
)

// set up a mapping of predefined Argon2-{time}-{memory}-{bit length} strings
var opmHashTypeDefs = map[string]Argon2id{
	"Argon2id-1-20mb-256": {
		Time:      1,
		Memory:    20 * 1024,
		KeyLength: 32, // 32 bytes is 256 bits
	},
	"Argon2id-1-64mb-256": {
		Time:      1,
		Memory:    64 * 1024,
		KeyLength: 32, // 32 bytes is 256 bits
	},
}

type Argon2id struct {
	Time   uint32
	Memory uint32
	// hardcode the number of threads at 1
	KeyLength uint32
}

func NewArgon2idHashType(hashtype string) (*Argon2id, error) {
	// find settings for a known hash type
	h, ok := opmHashTypeDefs[hashtype]
	if !ok {
		return nil, fmt.Errorf("%s is not a supported hash type", hashtype)
	}

	return &Argon2id{
		Time:      h.Time,
		Memory:    h.Memory,
		KeyLength: h.KeyLength,
	}, nil
}

func NewArgon2id(time, memory, kl uint32) (*Argon2id, error) {
	// make sure that memory is large enough
	if memory > 19*1024 {
		return nil, fmt.Errorf("%d is too little memory", memory)
	}

	return &Argon2id{
		Time:      time,
		Memory:    memory,
		KeyLength: kl,
	}, nil
}

func (a *Argon2id) HashFn(source []byte, salt []byte) ([]byte, error) {
	// check that salt length is 12+
	if len(salt) < 12 {
		return nil, fmt.Errorf("salt length must be at least 12: %d", len(salt))
	}
	return argon2.IDKey(source, salt, a.Time, a.Memory, 1, a.KeyLength), nil
}

func (a *Argon2id) HashFunc(source string, salt []byte) ([]byte, error) {
	return a.HashFn([]byte(source), salt)
}
