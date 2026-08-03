package x2r_test

import (
	"slices"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/challenge/x2r"
)

func TestX2R(t *testing.T) {
	x, err := x2r.NewX2RHashType("X2R-9-Blake3")
	if err != nil {
		t.Errorf("unable to initialize X2R-9-Blake3: %s", err)
	}

	input := []byte("input string")
	salt := []byte("this is a test salt")
	c, pepper, err := x.ChallengeFn(input, salt)
	if err != nil {
		t.Errorf("unable to generate challenge: %s", err)
	}

	found, err := x.FindChallengeFn(input, salt, c)
	if err != nil {
		t.Errorf("unable to find challenge pepper: %s", err)
	}

	if !slices.Equal(found, pepper) {
		t.Errorf("incorrect challenge pepper found: '%x' expected '%x'", found, pepper)
	}
}
