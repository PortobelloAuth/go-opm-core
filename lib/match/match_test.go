package match_test

import (
	"maps"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/lib/hashing"
	"github.com/PortobelloAuth/go-opm-core/lib/match"
	"github.com/PortobelloAuth/go-opm-core/lib/rules"
)

func TestMatch(t *testing.T) {
	noopHash, err := hashing.New("Noop")
	if err != nil {
		t.Fatalf("Unable to instatiate Noop (non-)hashing algorithm: %s", err)
	}

	raw := map[string]string{
		"given":  "Myra",
		"family": "Jones",
	}

	result, err := match.Match(
		&match.Demographics{
			Raw: raw,
		},
		&match.HashedDemographics{
			HashAlgorithm: noopHash.HashFunc,
			Salt:          []byte{},        // empty salt doesn't matter for Noop
			Hashes:        maps.Clone(raw), // copy the demographics
			Encode: func(b []byte) (string, error) {
				return string(b), nil
			},
			RuleSet: &rules.RuleSet{
				Rules: []rules.Rule{
					{Label: "fullname", Components: []rules.RuleDemographic{
						{Label: "given", Demographic: "given"},
						{Label: "family", Demographic: "family"},
					}},
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Unable to determine match status: %s", err)
	}

	if !result.Matched {
		t.Fatalf("Expected match status to be true: %v", result)
	}

	gtrue, gok := result.Detail["given"]
	ftrue, fok := result.Detail["family"]
	if !(gok && gtrue && fok && ftrue) {
		t.Fatalf("Expected match details for given and family to be true: %v", result)
	}
}
