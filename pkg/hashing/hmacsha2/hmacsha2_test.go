package hmacsha2_test

import (
	"encoding/hex"
	"slices"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/hashing/hmacsha2"
)

func TestHMACSHA2HashFn(t *testing.T) {
	expected, err := hex.DecodeString("76c48594567bb9ac90325a77c0fc0d7e1a10b948cac255543d22688e3aee4d92")
	if err != nil {
		t.Errorf("Error preparing to test NewSHA2(): %s", err)
	}

	h, err := hmacsha2.NewHMACSHA2("256")
	if err != nil {
		t.Errorf("Error creating NewSHA2(): %s", err)
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
