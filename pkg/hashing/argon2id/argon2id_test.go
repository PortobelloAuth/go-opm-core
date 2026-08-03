package argon2id_test

import (
	"encoding/hex"
	"slices"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/hashing/argon2id"
)

func TestArgon2idHashFn(t *testing.T) {
	expected, err := hex.DecodeString("18f6e6010a7db254b806256f0f94b57a38fe27376087e18a903a7dd25b035e7e")
	if err != nil {
		t.Errorf("Error preparing to test NewPBKDF2(): %s", err)
	}
	p, err := argon2id.NewArgon2id(1, 19*1024, 32)
	if err != nil {
		t.Errorf("Error creating NewArgon2id(): %s", err)
	}
	input := []byte("input string")
	salt := []byte("salt for testing")
	output, err := p.HashFn(input, salt)
	if err != nil {
		t.Errorf("Error generating hash with PBKDF2.HashFn(): %s", err)
	}

	if !slices.Equal(expected, output) {
		t.Errorf("HashFn output byte slice '%x' does not match expected slice %x", output, expected)
	}

	output2, err := p.HashFunc(string(input), salt)
	if err != nil {
		t.Errorf("Error generating hash with PBKDF2.HashFunc(): %s", err)
	}

	if !slices.Equal(expected, output2) {
		t.Errorf("HashFunc output byte slice '%x' does not match expected slice %x", output2, expected)
	}
}
