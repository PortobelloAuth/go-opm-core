package sha3_test

import (
	"encoding/hex"
	"slices"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/hashing/sha3"
)

func TestSHA3HashFn(t *testing.T) {
	expected, err := hex.DecodeString("90d527b23b431ee79c5bf1a8cf9bea1157f840c33a3d9ccdedf70327bf8ba902")
	if err != nil {
		t.Errorf("Error preparing to test NewSHA3(): %s", err)
	}

	h, err := sha3.NewSHA3("256")
	if err != nil {
		t.Errorf("Error creating NewSHA3(): %s", err)
	}

	input := []byte("input string")
	salt := []byte("salt for testing")
	output, err := h.HashFn(input, salt)
	if err != nil {
		t.Errorf("Error generating hash with sha3.HashFn(): %s", err)
	}

	if !slices.Equal(expected, output) {
		t.Errorf("HashFn output byte slice '%x' does not match expected slice %x", output, expected)
	}

	output2, err := h.HashFunc(string(input), salt)
	if err != nil {
		t.Errorf("Error generating hash with sha3.HashFunc(): %s", err)
	}

	if !slices.Equal(expected, output2) {
		t.Errorf("HashFunc output byte slice '%x' does not match expected slice %x", output2, expected)
	}
}
