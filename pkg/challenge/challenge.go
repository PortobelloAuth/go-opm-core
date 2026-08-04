package challenge

import (
	"fmt"
	"strings"

	"github.com/PortobelloAuth/go-opm-core/pkg/challenge/x2r"
)

type ChallengeConstructor func(specifier string) (Challenge, error)

var challengeRegistery map[string]ChallengeConstructor
var newCalled = false

func init() {
	challengeRegistery = make(map[string]ChallengeConstructor)

	// register the built-in Challenges
	builtins := []struct {
		K string
		N func(hashtype string) (Challenge, error)
	}{
		{"X2R", func(hashtype string) (Challenge, error) { return x2r.NewX2RHashType(hashtype) }},
	}

	for _, b := range builtins {
		challengeRegistery[b.K] = b.N
	}
}

func Register(hashtype string, fn ChallengeConstructor) error {
	// FIXME?: don't allow this to be called after New()? It makes testing annoying.
	if newCalled {
		return fmt.Errorf("Register() called after New()")
	}

	// FIXME?: don't allow overwriting builtins?
	challengeRegistery[hashtype] = fn

	return nil
}

type Challenge interface {
	ChallengeFn(source []byte, salt []byte) ([]byte, []byte, error)
	FindChallengeFn(source []byte, salt []byte, target []byte) ([]byte, error)
}

func New(hashtype string) (Challenge, error) {
	newCalled = true

	// split hashtype at the first hyphen
	prefix, _, found := strings.Cut(hashtype, "-")
	if found {
		fn, ok := challengeRegistery[prefix]
		if ok {
			return fn(hashtype)
		}
	}

	fn, ok := challengeRegistery[hashtype]
	if ok {
		return fn(hashtype)
	}

	return nil, fmt.Errorf("%s is not a registered challenge type", hashtype)
}
