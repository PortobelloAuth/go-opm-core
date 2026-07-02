package x2r

import (
	"encoding/binary"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/PortobelloAuth/go-opm-core/lib/hashing"
	"github.com/PortobelloAuth/go-opm-core/lib/randbytes"
)

// "Times 2 with Random bits"
// X2R takes another hashing function and runs it twice. First, it runs the basic
//
//	`base := HashFn(source, salt)`
//
// Second, it runs the hash function again with
//
//	`HashFn(append(base, randomBits...), salt)`
//
// where random bits is between 8 and 16 bits (at this time.) These bits are not
// disclosed to the reciever. This creates a challenge where the reciever must know
// the source material and discover the random bits through testing. The random bits
// can then be used to prove to the generator that the reciever solved the challenge.
type X2R struct {
	Hash    hashing.Hash
	RLength uint32
}

func NewX2RHashType(hashtype string) (*X2R, error) {
	// trim the XR2 prefix and et our hash function
	parts := strings.SplitN(hashtype, "-", 3)
	if parts[0] != "X2R" {
		return nil, fmt.Errorf("%s is not a known X2R hash type", hashtype)
	}

	// get rlen from the hash type
	rlen, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("Unable to recognize the number of random bits to include: %v", err)
	}

	subtype := parts[2]
	h, err := hashing.New(subtype)
	if err != nil {
		return nil, err
	}

	return &X2R{
		Hash:    h,
		RLength: uint32(rlen),
	}, nil
}

func (x *X2R) ChallengeFn(source []byte, salt []byte) ([]byte, []byte, error) {
	// check that salt length is 12+
	if len(salt) < 12 {
		return nil, nil, fmt.Errorf("salt length must be at least 12: %d", len(salt))
	}

	rbits, err := randbytes.RandomBits(x.RLength)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating random bits: %w", err)
	}

	h, err := x.hashChallenge(source, salt, rbits)
	if err != nil {
		return nil, nil, err
	}

	return h, rbits, nil
}

func (x *X2R) hashChallenge(source []byte, salt []byte, pepper []byte) ([]byte, error) {
	base, err := x.Hash.HashFn(source, salt)
	if err != nil {
		return nil, fmt.Errorf("error generating base hash: %w", err)
	}

	l2 := append(base, pepper...)
	h, err := x.Hash.HashFn(l2, salt)
	if err != nil {
		return nil, fmt.Errorf("error generating layer 2 hash: %w", err)
	}

	return h, nil
}

func (x *X2R) FindChallengeFn(source []byte, salt []byte, target []byte) ([]byte, error) {
	// check that salt length is 12+
	if len(salt) < 12 {
		return nil, fmt.Errorf("salt length must be at least 12: %d", len(salt))
	}

	// Loop over possible values of the random bits until we find a matching hash
	// NOTE: if the random bits used in the challenge remotely approaches 32 bits
	// this would take for...ev...er!
	buf := make([]byte, 4)
	max := uint32(1) << x.RLength
	for i := range max {
		binary.LittleEndian.PutUint32(buf, i)
		r, err := randbytes.MaskToBits(buf, x.RLength)
		if err != nil {
			return nil, err
		}
		h, err := x.hashChallenge(source, salt, r)
		if err != nil {
			return nil, err
		}
		if slices.Equal(h, target) {
			return r, nil
		}
	}

	return nil, fmt.Errorf("target hash not found")
}
