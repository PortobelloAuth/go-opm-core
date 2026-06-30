package randbytes_test

import (
	"encoding/binary"
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

func TestRandomBits(t *testing.T) {
	allZero := true
	// looping 16 times should be enough that at least 1 of our randomly generated
	// last bytes would have contained a high order bit before masking.
	for i := 0; i < 16; i++ {
		r, err := randbytes.RandomBits(12)
		if err != nil {
			t.Errorf("error generating random 12 bits: %s", err)
		}

		if len(r) != 2 {
			t.Errorf("wrong number of bytes for 12 random bits generated (should be 2): %d %0b", len(r), r)
		}

		lastbyte := r[len(r)-1]
		zero := lastbyte >> 4
		if zero != 0x00 {
			t.Errorf("unexpected set high-order bits found in last byte: %0b", lastbyte)
		}

		val := binary.BigEndian.Uint16(r)
		if val > 0 {
			// at least one of our 16 runs should not result in a zero value
			allZero = false
		}
	}

	if allZero {
		t.Errorf("all calls to randomBits resulted in a zero value")
	}
}
