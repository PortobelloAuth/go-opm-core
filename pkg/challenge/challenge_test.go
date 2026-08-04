package challenge_test

import (
	"slices"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/challenge"
)

func TestX2R(t *testing.T) {
	x, err := challenge.New("X2R-9-Blake3")
	if err != nil {
		t.Fatalf("unable to initialize X2R-9-Blake3: %s", err)
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
