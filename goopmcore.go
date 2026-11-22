package goopmcore

import (
	"fmt"

	"github.com/PortobelloAuth/go-opm-core/pkg/match"
	"github.com/PortobelloAuth/go-opm-core/pkg/match/matcher/plaintextcomponent"
	"github.com/PortobelloAuth/go-opm-core/pkg/rules"
)

// PlaintextMatch is a non-privacy preserving matching function to be used when both sets of identifying information
// are known locally. It is implemented using Match(), with a no-op hashing function, an empty salt, an encode
// function that converts the resulting byte slices back to strings, and the provided rule set.
// NOTE: this API explicitly ignores any normalization and hashing errors that may occur
// so that it can return a simple boolean value
func PlaintextMatch(check *match.Demographics, against *match.Demographics, withRules *rules.RuleSet) (bool, error) {
	plaintextmatcher, err := plaintextcomponent.New(withRules)
	if err != nil {
		return false, fmt.Errorf("Unable to instatiate plaintext Matcher: %w", err)
	}

	comp, _, err := plaintextmatcher.HashDemographics(against)
	if err != nil {
		return false, err
	}

	return Match(
		check,
		comp,
	)
}

// Match is the primary matching function of this library. It technically is a simple wrapper around
// match.Match().
// NOTE: this API explicitly ignores any normalization and hashing errors that occur
// so that it can return a simple boolean value. Callers interested in these errors should use
// match.Match() instead.
func Match(check *match.Demographics, against *match.HashedDemographics) (bool, error) {
	result, err := match.Match(
		check,
		against,
	)
	if err != nil {
		return false, err
	}

	return result.Matched, nil
}
