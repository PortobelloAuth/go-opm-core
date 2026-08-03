package sha2_test

import (
	"encoding/hex"
	"slices"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/hashing/sha2"
)

func TestSHA2_256HashFn(t *testing.T) {
	expected, err := hex.DecodeString("079c81678afd72ade9ab10f8bf130440644c8a7af7a506e0842d10ddae2b83df")
	if err != nil {
		t.Errorf("Error preparing to test NewPBKDF2(): %s", err)
	}

	h, err := sha2.NewSHA2("256")
	if err != nil {
		t.Errorf("Error creating NewSHA3(): %s", err)
	}

	input := []byte("input string")
	salt := []byte("salt for testing")
	output, err := h.HashFn(input, salt)
	if err != nil {
		t.Errorf("Error generating hash with sha2_256.HashFn(): %s", err)
	}

	if !slices.Equal(expected, output) {
		t.Errorf("HashFn output byte slice '%x' does not match expected slice %x", output, expected)
	}

	output2, err := h.HashFunc(string(input), salt)
	if err != nil {
		t.Errorf("Error generating hash with sha2_256.HashFunc(): %s", err)
	}

	if !slices.Equal(expected, output2) {
		t.Errorf("HashFunc output byte slice '%x' does not match expected slice %x", output2, expected)
	}
}
