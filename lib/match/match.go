package match

import (
	"fmt"

	"github.com/PortobelloAuth/go-opm-core/lib/hashing"
	"github.com/PortobelloAuth/go-opm-core/lib/rules"
)

type MatchResult struct {
	Matched bool
	// Detail contains information about which specifc demographics were and were not matched
	Detail map[string]bool
	// For now, MatchResult is a deterministic true/false result
	// In the future it may contain the threshold used in a binary statistical model along with a
	// confidence interval and other metadata in order to accomodate rules based on such models
}

type Demographics struct {
	Raw map[string]string
}

type HashedDemographics struct {
	HashAlgorithm hashing.HashFunc
	RuleSet       rules.RuleSet
	Hashes        map[string]string
}

// NewHashedDemographics creates HashedDemographics from Demographics by applying the RuleSet Normalizations
// and the supplied hashing function. Only demographics identified in the RuleSet are set in the resulting
// HashedDemographics Hashes map.
func NewHashedDemographics(raw *Demographics, hashFn hashing.HashFunc, r rules.RuleSet) (*HashedDemographics, error) {
	return nil, fmt.Errorf("Not implemented")
}

// Match takes a raw set of demographics and a hashed set of demographics, which contains a RuleSet, a
// hashing function, and the expected hashed demographic values. It processes the raw demographics
// using the RuleSet Normalizations, applies the hashing function, and compares the resulting values
// to the hashed set to see if any of the RuleSet's Rules match.
func Match(raw *Demographics, comp *HashedDemographics) (*MatchResult, error) {
	// Failed to match
	return &MatchResult{
		Matched: false,
		Detail:  make(map[string]bool),
	}, nil
}
