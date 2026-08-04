package hashing_test

import (
	"encoding/hex"
	"slices"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/hashing"
)

func TestBlake3HashFn(t *testing.T) {
	expected, err := hex.DecodeString("c6d25f73c520a48e5c193586a2265574d2278be2c16f6a8d91c8f087e4ba1b21")
	if err != nil {
		t.Errorf("Error preparing to test New(\"Blake3\"): %s", err)
	}

	h, err := hashing.New("Blake3")
	if err != nil {
		t.Errorf("Error creating New(\"Blake3\"): %s", err)
	}

	input := []byte("input string")
	salt := []byte("salt for testing")
	output, err := h.HashFn(input, salt)
	if err != nil {
		t.Errorf("Error generating hash with blake3.HashFn(): %s", err)
	}

	if !slices.Equal(expected, output) {
		t.Errorf("HashFn output byte slice '%x' does not match expected slice %x", output, expected)
	}

	output2, err := h.HashFunc(string(input), salt)
	if err != nil {
		t.Errorf("Error generating hash with blake3.HashFunc(): %s", err)
	}

	if !slices.Equal(expected, output2) {
		t.Errorf("HashFunc output byte slice '%x' does not match expected slice %x", output2, expected)
	}
}

func TestNoopHashFn(t *testing.T) {
	h, err := hashing.New("Noop")
	if err != nil {
		t.Errorf("Error creating hashing.New(\"Noop\"): %s", err)
	}

	input := []byte("input string")
	salt := []byte("test salt")
	output, err := h.HashFn(input, salt)
	if err != nil {
		t.Errorf("Error running noop.HashFn(): %s", err)
	}

	if !slices.Equal(input, output) {
		t.Errorf("Input byte slice '%s' does not match output slice '%s'", input, output)
	}
}
