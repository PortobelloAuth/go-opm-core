package componentmatcher

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
	"github.com/PortobelloAuth/go-opm-core/pkg/hashing"
	"github.com/PortobelloAuth/go-opm-core/pkg/match"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization"
	"github.com/PortobelloAuth/go-opm-core/pkg/rules"
)

// Matcher is a match.Matcher with the various fields that used to be in HashedDemographics
// to try to make the Match algortihm sufficiently general. It implements the rule-component
// matching algorithm that was previously built in to HashedDemographics. That matching
// algorithm was premised on the idea that the hashing function (and therefore the salt and
// encoding function) was the only real difference between matching plaintext and matching
// hashes. As we began to implement fuzzy matching it became apparent that some hashing
// strategies are deeply interconnected with the matching algorithm they require.
type Matcher struct {
	HashAlgorithm    hashing.HashFunc
	Salt             []byte
	Encode           func([]byte) (string, error)
	Rules            *rules.RuleSet
	flattenedRuleSet rules.FlattenedRuleSet
}

func New(rs *rules.RuleSet, h hashing.HashFunc, s []byte, enc func([]byte) (string, error)) (*Matcher, error) {
	flatrs, err := rs.Flatten()
	if err != nil {
		return nil, fmt.Errorf("unable to flatten RuleSet: %w", err)
	}

	return &Matcher{
		HashAlgorithm:    h,
		Salt:             s,
		Encode:           enc,
		Rules:            rs,
		flattenedRuleSet: flatrs,
	}, nil
}

func ruleComponentKey(label string, key string) string {
	return fmt.Sprintf("%s|%s", label, key)
}

// HashDemographics creates HashedDemographics from Demographics by applying the Matcher's RuleSet
// Normalizations and hashing function.. Only demographics identified in the RuleSet are set in the
// resulting HashedDemographics Hashes map.
func (m *Matcher) HashDemographics(raw *match.Demographics) (*match.HashedDemographics, error, error) {
	// walk the RuleSet tree collecting the demographic keys and Normlizations that need to be applied. Prefix
	// the Rule.Label to demographics so that hashes for different rules are mapped appropriately.

	flatrules := m.flattenedRuleSet

	hashErrors := make([]error, 0)
	hashes := maps.Collect(func(yield func(string, []string) bool) {
		for _, flatrule := range flatrules {
			for _, rd := range flatrule.Components {
				v, ok := raw.Raw[rd.Key]
				if ok {
					normv := v
					if rd.Normalizations != nil {
						var norm = normalization.Normalization(rd.Normalizations)
						with, ok := norm.(normalization.ContextualNormalization)
						if ok {
							norm = with.WithContext(func(s string) (string, bool) {
								v, ok := raw.Raw[s]
								return v, ok
							})
						}
						var err error
						normv, err = norm.Apply(v)
						if err != nil {
							hashErrors = append(hashErrors, fmt.Errorf("Error normalizing %q for %q: %w", rd.Key, flatrule.Label, err))
							continue
						}
					}

					altFn := demographicalts.IdentityDemographicAlternativesFn
					if rd.Alternatives != nil {
						altFn = rd.Alternatives
					}
					alts := altFn(normv, rd.Key, rd.Limit)
					althashes := slices.Collect(func(yield func(string) bool) {
						for _, a := range alts {
							if a.Text == "" {
								// empty normalized values MUST NOT be compared. DO NOT hash them.
								// demographicalts as implemented here depend on the normalized value.
								// TODO: Other alternative values (like nicknames and historical names)
								// may require their own implementation or an adjustment to this logic.
								// This check is implemented here instead of before the altFn is called
								// to leave room for possible implementations.
								continue
							}
							althashbytes, err := m.HashAlgorithm(a.Text, m.Salt)
							if err != nil {
								hashErrors = append(hashErrors, fmt.Errorf("Error hashing %q for %q: %w", rd.Key, flatrule.Label, err))
								continue
							}
							althash, err := m.Encode(althashbytes)
							if err != nil {
								hashErrors = append(hashErrors, fmt.Errorf("Error encoding %q hash for %q: %w", rd.Key, flatrule.Label, err))
								continue
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

	return &match.HashedDemographics{
		Matcher: m,
		Hashes:  hashes,
	}, errors.Join(hashErrors...), nil
}

func (m *Matcher) Match(calc, comp *match.HashedDemographics) (*match.MatchResult, error) {
	// check the resulting hashes maps, which should be keyed by the rule labels for the rules that
	// generate them, to see which rules actually matched. The values in the hash maps are lists of
	// hashes to support demographic alternatives (nick names, fuzzy matching, etc.). Only one hash
	// for a Rule needs to match for the Rule to match. Only one Rule for a RuleSet needs to match.
	detail := make(map[string]bool)
	matched := false

	// iterate the flattened rule set for rules to check
	for _, flatrule := range m.flattenedRuleSet {
		matchedRule := 0

		// check each component of the flat rule to see if it has a matching hash
		for _, rd := range flatrule.Components {
			hashkey := ruleComponentKey(flatrule.Label, rd.Key)

			detail[hashkey] = false
			have, ok := calc.Hashes[hashkey]
			if !ok {
				// we don't have the calculated hash; check the next rule
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
					if havehash == comphash {
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

	return &match.MatchResult{
		Matched: matched,
		Detail:  detail,
	}, nil
}
