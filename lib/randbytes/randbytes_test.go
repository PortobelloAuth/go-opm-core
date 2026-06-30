package randbytes_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/lib/randbytes"
)

func TestRandomBytes(t *testing.T) {
	r, err := randbytes.RandomBytes(16)
	if err != nil {
		t.Errorf("error generating random 16 bytes: %s", err)
	}

	if len(r) != 16 {
		t.Errorf("wrong number of random bytes generated: %d", len(r))
	}

	// NOTE: No test is done to check that the implementation is using a suitable random number generator
}

func TestGenerateSalt(t *testing.T) {
	r1, err := randbytes.GenerateSalt(2)
	if r1 != nil {
		t.Errorf("unexpected byte slice returned from call to GenerateSalt() with too small a length")
	}
	if err == nil {
		t.Errorf("no error returned from call to GenerateSalt() with too small a length")
	}

	r2, err := randbytes.GenerateSalt(20)
	if err != nil {
		t.Errorf("error returned from call to GenerateSalt(): %s", err)
	}
	if r2 == nil {
		t.Errorf("nil byte slice returned from call to GenerateSalt()")
	}
}
