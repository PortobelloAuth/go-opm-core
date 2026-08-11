package match_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
	_ "github.com/PortobelloAuth/go-opm-core/pkg/demographicalts/diacriticalts"
	"github.com/PortobelloAuth/go-opm-core/pkg/hashing"
	"github.com/PortobelloAuth/go-opm-core/pkg/match"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization"
	"github.com/PortobelloAuth/go-opm-core/pkg/rules"
)

var testcases = []struct {
	Raw    map[string]string
	Hashes map[string][]string
}{
	{
		Raw: map[string]string{
			"given":  "Maria",
			"family": "García",
		},
		Hashes: map[string][]string{
			"fullname|given":  {"maria"},
			"fullname|family": {"garcía", "garcia"},
		},
	},
	{
		Raw: map[string]string{
			"given":  "José",
			"family": "López",
		},
		Hashes: map[string][]string{
			"fullname|given":  {"josé"},
			"fullname|family": {"lópez", "lopez"},
		},
	},
	{
		Raw: map[string]string{
			"given":  "Myra",
			"family": "Jones",
		},
		Hashes: map[string][]string{
			"fullname|given":  {"myra"},
			"fullname|family": {"jones"},
		},
	},
}

func TestMatch(t *testing.T) {
	noopHash, err := hashing.New("Noop")
	if err != nil {
		t.Fatalf("Unable to instatiate Noop (non-)hashing algorithm: %s", err)
	}

	noopSalt := []byte{} // empty salt doesn't matter for Noop

	lower, err := normalization.New("lowercase")
	if err != nil {
		t.Fatalf("Failed to initialize %s normalization", "lowercase")
	}

	diacritics, err := demographicalts.New("diacritics")
	if err != nil {
		t.Fatalf("Failed to initialize %s alternatives function", "diacritics")
	}

	testRules := &rules.RuleSet{
		Normalizations: &normalization.NormalizationSequence{Seq: []normalization.Normalization{lower}},
		Rules: []rules.Rule{
			{Label: "fullname", Components: []rules.RuleDemographic{
				{Label: "given", Key: "given"},
				{Label: "family", Key: "family", Alternatives: diacritics},
			}},
		},
	}

	for _, tc := range testcases {
		rawDemographics := &match.Demographics{
			Raw: tc.Raw,
		}

		result, err := match.Match(
			rawDemographics,
			&match.HashedDemographics{
				HashAlgorithm: noopHash.HashFunc,
				Salt:          noopSalt,
				Hashes:        tc.Hashes,
				Encode:        match.HashToString,
				RuleSet:       testRules,
			},
		)
		if err != nil {
			t.Fatalf("Unable to determine match status: %s", err)
		}

		if !result.Matched {
			t.Fatalf("Expected match status to be true: %v", result)
		}

		for hkey, _ := range tc.Hashes {
			d, ok := result.Detail[hkey]
			if !(ok && d) {
				t.Fatalf("Expected match details for %s to be true: %v", hkey, result)
			}
		}
	}
}
