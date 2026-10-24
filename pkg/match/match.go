package match

import (
	"fmt"
	"maps"
	"slices"

	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
	"github.com/PortobelloAuth/go-opm-core/pkg/hashing"
	"github.com/PortobelloAuth/go-opm-core/pkg/rules"
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
	HashAlgorithm    hashing.HashFunc
	Salt             []byte
	RuleSet          *rules.RuleSet
	flattenedRuleSet rules.FlattenedRuleSet
	// Hahses contains a map of normalized and hashed demographic alternatives, keyed by the flattened
	// rule label that produced it.
	Hashes map[string][]string
	Encode func([]byte) (string, error)
}

func ruleComponentKey(label string, key string) string {
	return fmt.Sprintf("%s|%s", label, key)
}

// NewHashedDemographics creates HashedDemographics from Demographics by applying the RuleSet Normalizations
// and the supplied hashing function. Only demographics identified in the RuleSet are set in the resulting
// HashedDemographics Hashes map.
func NewHashedDemographics(raw *Demographics, hashFn hashing.HashFunc, salt []byte, rset *rules.RuleSet, encode func([]byte) (string, error)) (*HashedDemographics, error) {
	// walk the RuleSet tree collecting the demographic keys and Normlizations that need to be applied. Prefix
	// the Rule.Label to demographics so that hashes for different rules are mapped appropriately.

	flatrs, err := rset.Flatten()
	if err != nil {
		return nil, fmt.Errorf("Unable to apply RuleSet: %w", err)
	}

	hashErrors := make([]error, 0)
	hashes := maps.Collect(func(yield func(string, []string) bool) {
		for _, flatrule := range flatrs {
			for _, rd := range flatrule.Components {
				v, ok := raw.Raw[rd.Key]
				if ok {
					normv := v
					if rd.Normalizations != nil {
						var err error
						normv, err = rd.Normalizations.Apply(v)
						if err != nil {
							hashErrors = append(hashErrors, fmt.Errorf("Error normalizing %s for %s: %w", rd.Key, flatrule.Label, err))
							return
						}
					}

					altFn := demographicalts.IdentityDemographicAlternativesFn
					if rd.Alternatives != nil {
						altFn = rd.Alternatives
					}
					alts := altFn(normv, rd.Key, rd.Limit)
					althashes := slices.Collect(func(yield func(string) bool) {
						for _, a := range alts {
							althashbytes, err := hashFn(a.Text, salt)
							if err != nil {
								hashErrors = append(hashErrors, fmt.Errorf("Error hashing %s for %s: %w", rd.Key, flatrule.Label, err))
								return
							}
							althash, err := encode(althashbytes)
							if err != nil {
								hashErrors = append(hashErrors, fmt.Errorf("Error encoding %s hash for %s: %w", rd.Key, flatrule.Label, err))
								return
							}
							if !yield(althash) {
								return
							}
						}
					})
					hashkey := ruleComponentKey(flatrule.Label, rd.Key)
					if !yield(hashkey, althashes) {
						return
					}
				}
			}
		}
	})

	if len(hashErrors) > 0 {
		return nil, hashErrors[0]
	}

	return &HashedDemographics{
		HashAlgorithm:    hashFn,
		Salt:             salt,
		RuleSet:          rset,
		flattenedRuleSet: flatrs,
		Hashes:           hashes,
		Encode:           encode,
	}, nil
}

// Match takes a raw set of demographics and a hashed set of demographics, which contains a RuleSet, a
// hashing function, and the expected hashed demographic values. It processes the raw demographics
// using the RuleSet Normalizations, applies the hashing function, and compares the resulting values
// to the hashed set to see if any of the RuleSet's Rules match.
func Match(raw *Demographics, comp *HashedDemographics) (*MatchResult, error) {
	calc, err := NewHashedDemographics(
		raw,
		comp.HashAlgorithm,
		comp.Salt,
		comp.RuleSet,
		comp.Encode,
	)
	if err != nil {
		return nil, fmt.Errorf("Error calculating matching hashes: %w", err)
	}

	return calc.Match(comp)
}

func (calc *HashedDemographics) Match(comp *HashedDemographics) (*MatchResult, error) {

	// check the resulting hashes maps, which should be keyed by the rule labels for the rules that
	// generate them, to see which rules actually matched. The values in the hash maps are lists of
	// hashes to support demographic alternatives (nick names, fuzzy matching, etc.). Only one hash
	// for a Rule needs to match for the Rule to match. Only one Rule for a RuleSet needs to match.
	detail := make(map[string]bool)
	matched := false

	// iterate the flattened rule set for rules to check
	for _, flatrule := range calc.flattenedRuleSet {
		matchedRule := 0

		// check each component of the flat rule to see if it has a matching hash
		for _, rd := range flatrule.Components {
			hashkey := ruleComponentKey(flatrule.Label, rd.Key)

			detail[hashkey] = false
			have, ok := calc.Hashes[hashkey]
			if !ok {
				// we don't have the calculated hash (!?!?); check the next rule
				continue
			}
			check, ok := comp.Hashes[hashkey]
			if !ok {
				// we don't have the comparison hash; check the next rule
				continue
			}
			found := false
			for _, havehash := range have {

				for _, comphash := range check {
					// only one of the hash alternatives needs to match for each component
					if ok && havehash == comphash {
						found = true
					}
				}
			}
			if found {
				detail[hashkey] = true
				matchedRule += 1
			}
		}

		if matchedRule == len(flatrule.Components) {
			matched = true
		}
	}

	return &MatchResult{
		Matched: matched,
		// FIXME: fill out the Detail object!!
		Detail: detail,
	}, nil
}

func HashToString(b []byte) (string, error) {
	return string(b), nil
}
