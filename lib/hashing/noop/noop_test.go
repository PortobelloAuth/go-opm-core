package noop_test

import (
	"slices"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/lib/hashing/noop"
)

func TestNoopHashFn(t *testing.T) {
	h, err := noop.New()
	if err != nil {
		t.Errorf("Error creating noop.New(): %s", err)
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

func TestNoopHashFunc(t *testing.T) {
	h, err := noop.New()
	if err != nil {
		t.Errorf("Error creating noop.New(): %s", err)
	}

	input := "input string"
	salt := []byte("test salt")
	output, err := h.HashFunc(input, salt)
	if err != nil {
		t.Errorf("Error running noop.HashFunc(): %s", err)
	}

	if input != string(output) {
		t.Errorf("Input string '%s' does not match output string '%s'", input, output)
	}
}
