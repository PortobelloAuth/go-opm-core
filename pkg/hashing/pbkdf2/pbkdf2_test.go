package pbkdf2_test

import (
	"encoding/hex"
	"slices"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/hashing/pbkdf2"
)

func TestPBKDF2HashFn(t *testing.T) {
	expected, err := hex.DecodeString("467bf40c00a4b4b9b907044ed11ab1d3876b93f9189dee5c228b0b55cd872552")
	if err != nil {
		t.Errorf("Error preparing to test NewPBKDF2(): %s", err)
	}
	p, err := pbkdf2.NewPBKDF2("SHA2_256", 10000, 32)
	if err != nil {
		t.Errorf("Error creating NewPBKDF2(): %s", err)
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
