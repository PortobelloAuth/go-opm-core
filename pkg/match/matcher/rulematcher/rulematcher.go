package rulematcher

import (
	"errors"
	"fmt"
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
	HashLabels       bool
	Salt             []byte
	Encode           func([]byte) (string, error)
	Rules            *rules.RuleSet
	flattenedRuleSet rules.FlattenedRuleSet
}

func New(rs *rules.RuleSet, hashlabels bool, h hashing.HashFunc, s []byte, enc func([]byte) (string, error)) *Matcher {
	return &Matcher{
		HashAlgorithm: h,
		HashLabels:    hashlabels,
		Salt:          s,
		Encode:        enc,
		Rules:         rs,
	}
}

type RuleValuePartial struct {
	Label string
	Value string
}

func (m *Matcher) hashAndEncode(v string) (string, error) {
	hashbytes, err := m.HashAlgorithm(v, m.Salt)
	if err != nil {
		return "", fmt.Errorf("hashing error: %w", err)
	}
	hashv, err := m.Encode(hashbytes)
	if err != nil {
		return "", fmt.Errorf("encoding error: %w", err)
	}
	return hashv, nil
}

// HashDemographics creates HashedDemographics from Demographics by applying the Matcher's RuleSet
// Normalizations and hashing function.. Only demographics identified in the RuleSet are set in the
// resulting HashedDemographics Hashes map.
func (m *Matcher) HashDemographics(raw *match.Demographics) (*match.HashedDemographics, error, error) {
	// walk the RuleSet tree collecting the demographic keys and Normlizations that need to be applied. Prefix
	// the Rule.Label to demographics so that hashes for different rules are mapped appropriately.

	flatrules := m.flattenedRuleSet
	if flatrules == nil {
		flatrs, err := m.Rules.Flatten()
		if err != nil {
			return nil, nil, fmt.Errorf("Unable to apply RuleSet: %w", err)
		}
		m.flattenedRuleSet = flatrs
		flatrules = flatrs
	}

	hashErrors := make([]error, 0)
	hashes := map[string][]string{}
	labelhashes := map[string]string{}
	for _, flatrule := range flatrules {
		// a list of the rule inputs so far
		partials := make([]RuleValuePartial, 0)

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
				// For each alt and each partial (if there are any) add a new partial
				newpartials := slices.Collect(func(yield func(RuleValuePartial) bool) {
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

						if len(partials) == 0 {
							if !yield(RuleValuePartial{
								Label: fmt.Sprintf("%s:%s", flatrule.Label, a.Label),
								Value: a.Text,
							}) {
								return
							}
						} else {
							for _, p := range partials {
								if !yield(RuleValuePartial{
									Label: fmt.Sprintf("%s-%s", p.Label, a.Label),
									Value: fmt.Sprintf("%s|%s", p.Value, a.Text),
								}) {
									return
								}
							}
						}
					}
				})
				partials = newpartials
			}
		}

		for _, p := range partials {
			hashv, err := m.hashAndEncode(p.Value)
			if err != nil {
				hashErrors = append(hashErrors, fmt.Errorf("Error hashing %q for %q: %w", p.Label, flatrule.Label, err))
				continue
			}

			hashkey := p.Label
			if m.HashLabels {
				hashkey, err = m.hashAndEncode(p.Label)
				if err != nil {
					hashErrors = append(hashErrors, fmt.Errorf("Error hashing key %q for %q: %w", p.Label, flatrule.Label, err))
					continue
				}
			}
			labelhashes[hashkey] = p.Label

			existing, ok := hashes[hashkey]
			if !ok {
				existing = make([]string, 0)
			}
			existing = append(existing, hashv)
			hashes[hashkey] = existing
		}
	}

	return &match.HashedDemographics{
		Matcher: m,
		Hashes:  hashes,
		Aux:     labelhashes,
	}, errors.Join(hashErrors...), nil
}

func (m *Matcher) Match(calc, comp *match.HashedDemographics) (*match.MatchResult, error) {
	fmt.Printf("comparing\n\tcalc: %v\n\tcomp: %v\n", calc, comp)
	// check the resulting hashes maps, which should be keyed by the rule labels for the rules that
	// generate them, to see which rules actually matched. The values in the hash maps are lists of
	// hashes to support demographic alternatives (nick names, fuzzy matching, etc.). Only one hash
	// for a Rule needs to match for the Rule to match. Only one Rule for a RuleSet needs to match.
	detail := make(map[string]bool)
	matched := false
	labelhashes, ok := calc.Aux.(map[string]string)
	if !ok {
		labelhashes = map[string]string{}
	}

	// the rule label and its hash must exist in both, so we can iterate just one and be just fine
	for labelhash, rulehashes := range calc.Hashes {
		comphashes, ok := comp.Hashes[labelhash]
		if !ok {
			continue
		}
		label, ok := labelhashes[labelhash]
		if !ok {
			label = labelhash
		}

		for _, rh := range rulehashes {
			if slices.Contains(comphashes, rh) {
				matched = true
				detail[label] = true
				if m.HashLabels {
					// put the hashed label in detail as well, in case it is needed
					detail[labelhash] = true
				}
			}
		}
	}

	return &match.MatchResult{
		Matched: matched,
		Detail:  detail,
	}, nil
}
