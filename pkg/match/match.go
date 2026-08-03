package match

import (
	"fmt"

	"github.com/PortobelloAuth/go-opm-core/pkg/hashing"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization"
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
	HashAlgorithm hashing.HashFunc
	Salt          []byte
	RuleSet       *rules.RuleSet
	Hashes        map[string]string
	Encode        func([]byte) (string, error)
}

// NewHashedDemographics creates HashedDemographics from Demographics by applying the RuleSet Normalizations
// and the supplied hashing function. Only demographics identified in the RuleSet are set in the resulting
// HashedDemographics Hashes map.
func NewHashedDemographics(raw *Demographics, hashFn hashing.HashFunc, salt []byte, rset *rules.RuleSet, encode func([]byte) (string, error)) (*HashedDemographics, error) {
	// walk the RuleSet tree collecting the demographic keys and Normlizations that need to be applied. Prefix
	// the Rule.Label to demographics so that hashes for different rules are mapped appropriately.
	hashes := make(map[string]string, 0)
	err := recursiveHashedDemographics(rset.Normalizations, raw, hashFn, salt, rset, encode, hashes)
	if err != nil {
		return nil, fmt.Errorf("Unable to apply RuleSet: %w", err)
	}

	return &HashedDemographics{
		HashAlgorithm: hashFn,
		Salt:          salt,
		RuleSet:       rset,
		Hashes:        hashes,
		Encode:        encode,
	}, nil
}

func recursiveHashedDemographics(parentnorms *normalization.NormalizationSequence, raw *Demographics, hashFn hashing.HashFunc, salt []byte, rset *rules.RuleSet, encode func([]byte) (string, error), collected map[string]string) error {
	for _, rule := range rset.Rules {
		if rule.Set != nil {
			sequence := []normalization.Normalization{parentnorms, rule.Set.Normalizations}
			err := recursiveHashedDemographics(
				&normalization.NormalizationSequence{Seq: sequence},
				raw,
				hashFn,
				salt,
				rule.Set,
				encode,
				collected,
			)
			if err != nil {
				return err
			}
			// A Rule may only contain a RuleSet OR a Components, not both
			continue
		}
		// NOTE: (FIXME) avoid recomputing hashes for rules that are using the same demographic,
		// normalizations, and fuzzy matching logic.
		// TODO: all rule components: demographics, normalizations, alternative functions, etc. have (or can have)
		// a unique namespaced identifier (aka specifier) that exists as part of the Rule configuration DSL (which
		// is a subset of the NetworkSettings DSL.) A defined, ordered serialization of thes specifiers should
		// allow us to use them to build a key to look up the hashes we need. The original strawman spec's DSL
		// might be (mostly) sufficient.
		for _, d := range rule.Components {
			var err error
			// Run each component through d.Alternatives(), etc.
			fmt.Printf("Demographic: %v\n", d)
			rawvalue, ok := raw.Raw[d.Demographic]
			if !ok || len(rawvalue) < 1 {
				// skip
				continue
			}
			altFn := d.Alternatives
			if altFn == nil {
				altFn = rules.IdentityDemographicAlternativesFn
			}
			alts := altFn(rawvalue, d.Label, d.Limit)
			// Now properly populate the collected map
			for _, a := range alts {
				n := a.Text
				if parentnorms != nil {
					n, err = parentnorms.Apply(a.Text)
					if err != nil {
						return fmt.Errorf("failed to normalize text for %s: %w", a.Label, err)
					}
				}
				h, err := hashFn(n, salt)
				if err != nil {
					return fmt.Errorf("failed to generte hash for %s: %w", a.Label, err)
				}
				// encode hashes to string (in most cases base64url)
				hstr, err := encode(h)
				if err != nil {
					return fmt.Errorf("failed to encode hash for %s: %w", a.Label, err)
				}
				collected[a.Label] = hstr
			}
		}
	}
	return nil
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
		return nil, fmt.Errorf("Error during matching: %w", err)
	}

	fmt.Printf("calculated hashes: %v\n", calc)
	// TODO: check the Rules to see if the needed hashes match. The Hashes object should only be
	// one level deep, so, if each hash represents a full rule variant, we just need to search
	// the comp object and calc object for keys that match and have matching values.
	// FIXME: without the DSL to define hashes on a per variant basis, this check isn't quite valid!!
	matched := false
	for k, expected := range comp.Hashes {
		found, ok := calc.Hashes[k]
		if ok && found == expected {
			matched = true
		}
	}

	// Failed to match
	return &MatchResult{
		Matched: matched,
		// FIXME: fill out the Detail object!!
		Detail: make(map[string]bool),
	}, nil
}
