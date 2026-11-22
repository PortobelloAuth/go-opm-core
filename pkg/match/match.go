package match

import "fmt"

type MatchResult struct {
	Matched bool
	// Detail contains information about which specifc demographics were and were not matched
	Detail map[string]bool
	// HashError contains any normalization and hashing errors encountering during hashing.
	// These errors specifically DO NOT block matching, though they do prevent the use of any
	// Rules that rely on them. They are therefore returned as part of the MatchResult, not
	// as distinct error return value.
	HashError error
	// For now, MatchResult is a deterministic true/false result
	// In the future it may contain the threshold used in a binary statistical model along with a
	// confidence interval and other metadata in order to accomodate rules based on such models
}

type Demographics struct {
	Raw map[string]string
}

type HashedDemographics struct {
	// Hashes contains a map of normalized and hashed demographic alternatives, keyed by the flattened
	// rule label that produced it.
	Hashes  map[string][]string
	Matcher Matcher
	// Aux exists to allow a Matcher to optionally associate Matcher-specific data with the HashedDemographics
	// instance during the HashDemographics() call that will be needed later in the Match() call.
	Aux interface{}
}

type Matcher interface {
	// NOTE: HashDemographics returns the HashedDemographics and two potentially non-nil errors. The first
	// contains all of the normalization and hashing errors that occured while applying the RuleSet to the
	// raw Demographics. This error is informational and is added to the MatchResult by Match() as HashError. The
	// second error carries fatal errors indicating that the whole hash operation has failed. Callers must be
	// aware of both kinds of errors.
	HashDemographics(raw *Demographics) (*HashedDemographics, error, error)
	Match(calc, comp *HashedDemographics) (*MatchResult, error)
}

// Match takes a raw set of demographics and a hashed set of demographics, which contains a RuleSet, a
// hashing function, and the expected hashed demographic values. It processes the raw demographics
// using the RuleSet Normalizations, applies the hashing function, and compares the resulting values
// to the hashed set to see if any of the RuleSet's Rules match.
func Match(raw *Demographics, comp *HashedDemographics) (*MatchResult, error) {
	calc, hasherr, err := comp.Matcher.HashDemographics(raw)
	if err != nil {
		return nil, err
	}
	if calc == nil {
		return nil, fmt.Errorf("Unable to calculate matching hashes: %w", hasherr)
	}

	matched, err := comp.Matcher.Match(calc, comp)
	if err != nil {
		return nil, err
	}
	matched.HashError = hasherr

	return matched, nil
}
