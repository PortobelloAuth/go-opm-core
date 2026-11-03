package match_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
	_ "github.com/PortobelloAuth/go-opm-core/pkg/demographicalts/diacriticalts"
	"github.com/PortobelloAuth/go-opm-core/pkg/hashing"
	"github.com/PortobelloAuth/go-opm-core/pkg/hashing/noop"
	"github.com/PortobelloAuth/go-opm-core/pkg/match"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization"
	"github.com/PortobelloAuth/go-opm-core/pkg/rules"
)

func TestMatch(t *testing.T) {
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
			t.Fatalf("Unexpected error determining match status: %s", err)
		}

		if !result.Matched {
			t.Fatalf("Expected match status to be true: %v", result)
		}

		for hkey := range tc.Hashes {
			d, ok := result.Detail[hkey]
			if !(ok && d) {
				t.Fatalf("Expected match details for %s to be true: %v", hkey, result)
			}
		}
	}
}

// ErrorNormalizer is a no-op normalization function that errors when it
// sees the string "normal" so that we can test what Match() does when
// encountering a normalization error.
type ErrorNormalizer struct{}

func NewErrorNormalizer() (*ErrorNormalizer, error) {
	return &ErrorNormalizer{}, nil
}

var testNormalizationError = errors.New("test normalization error")

func (l *ErrorNormalizer) Apply(input string) (string, error) {
	if strings.Contains(input, "normal") {
		return "", testNormalizationError
	}
	return input, nil
}

// ErrorNoop is a no-op hashing function that errors when it sees the
// string "hash" so that we can test what Match() does when encountering
// a hashing error.
type ErrorNoop struct{}

func NewErrorNoop() (*ErrorNoop, error) {
	return &ErrorNoop{}, nil
}

func (n *ErrorNoop) HashFn(source []byte, salt []byte) ([]byte, error) {
	return n.HashFunc(string(source), salt)
}

var testHashError = errors.New("test hash error")

func (n *ErrorNoop) HashFunc(source string, salt []byte) ([]byte, error) {
	if strings.Contains(source, "hash") {
		return nil, testHashError
	}
	return []byte(source), nil
}

func TestMatchWithErrors(t *testing.T) {
	var testcases = []struct {
		Raw    map[string]string
		Hashes map[string][]string
		Result bool
		Error  error
	}{
		{
			Raw: map[string]string{
				"given":  "Maria",
				"family": "García",
				"error":  "",
			},
			Hashes: map[string][]string{
				"fullname|given":  {"maria"},
				"fullname|family": {"garcía", "garcia"},
			},
			Result: true,
			Error:  nil,
		},
		{
			Raw: map[string]string{
				"given":  "José",
				"family": "López",
				"error":  "normalization",
			},
			Hashes: map[string][]string{
				"fullname|given":  {"josé"},
				"fullname|family": {"lópez", "lopez"},
			},
			Result: true,
			Error:  testNormalizationError,
		},
		{
			Raw: map[string]string{
				"given":  "Myra",
				"family": "Jones",
				"error":  "hashing",
			},
			Hashes: map[string][]string{
				"fullname|given":  {"myra"},
				"fullname|family": {"jones"},
			},
			Result: true,
			Error:  testHashError,
		},
		{
			Raw: map[string]string{
				"given":  "Myra",
				"family": "Jones",
				"error":  "hashing",
			},
			Hashes: map[string][]string{
				"fullname|given":  {"adam"},
				"fullname|family": {"smith"},
			},
			Result: false,
			Error:  testHashError,
		},
	}

	noopHash, err := NewErrorNoop()
	if err != nil {
		t.Fatalf("Unable to instatiate ErrorNoop test (non-)hashing algorithm: %s", err)
	}

	noopSalt := []byte{} // empty salt doesn't matter for Noop

	lower, err := normalization.New("lowercase")
	if err != nil {
		t.Fatalf("Failed to initialize %s normalization", "lowercase")
	}

	errnorm, err := NewErrorNormalizer()
	if err != nil {
		t.Fatalf("Failed to initialize test erroring normalization")
	}

	diacritics, err := demographicalts.New("diacritics")
	if err != nil {
		t.Fatalf("Failed to initialize %s alternatives function", "diacritics")
	}

	testRules := &rules.RuleSet{
		Normalizations: &normalization.NormalizationSequence{Seq: []normalization.Normalization{errnorm, lower}},
		Rules: []rules.Rule{
			{Label: "fullname", Components: []rules.RuleDemographic{
				{Label: "given", Key: "given"},
				{Label: "family", Key: "family", Alternatives: diacritics},
			}},
			{Label: "other", Components: []rules.RuleDemographic{
				{Label: "error", Key: "error"},
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
			t.Fatalf("Unexpected returned error from Match() %q", err)
		}

		if !errors.Is(result.HashError, tc.Error) {
			t.Fatalf("Expected match result HashError %q to be %q", err, tc.Error)
		}

		if tc.Result != result.Matched {
			t.Fatalf("Expected match status to be %t: %v", tc.Result, result)
		}

		if tc.Result {
			for hkey := range tc.Hashes {
				d, ok := result.Detail[hkey]
				if !(ok && d) {
					t.Fatalf("Expected match details for %s to be true: %v", hkey, result)
				}
			}
		}
	}
}

// EmptyNormalizer is a normalization function turns any value under an .
type EmptyNormalizer struct{}

func NewEmptyNormalizer() (*EmptyNormalizer, error) {
	return &EmptyNormalizer{}, nil
}

func (l *EmptyNormalizer) Apply(input string) (string, error) {
	return "", nil
}

func TestMatchWithBlanks(t *testing.T) {
	var testcases = []struct {
		RawA map[string]string
		RawB map[string]string
	}{
		{
			RawA: map[string]string{
				"given":  "Maria",
				"family": "García",
				"empty":  "",
			},
			RawB: map[string]string{
				"given":  "maria",
				"family": "garcía",
			},
		},
		{
			RawA: map[string]string{
				"given":  "José",
				"family": "López",
				"empty":  "normalization",
			},
			RawB: map[string]string{
				"given":  "josé",
				"family": "lópez",
			},
		},
		{
			RawA: map[string]string{
				"given":  "Myra",
				"family": "Jones",
				"empty":  "",
			},
			RawB: map[string]string{
				"given":  "myra",
				"family": "jones",
			},
		},
		{
			RawA: map[string]string{
				"given":  "Myra",
				"family": "Jones",
				"empty":  "",
			},
			RawB: map[string]string{
				"given":  "adam",
				"family": "smith",
			},
		},
		{
			RawA: map[string]string{
				"empty": "",
			},
			RawB: map[string]string{},
		},
	}

	noopHash, err := noop.New()
	if err != nil {
		t.Fatalf("Unable to instatiate noop (non-)hashing algorithm: %s", err)
	}

	noopSalt := []byte{} // empty salt doesn't matter for Noop

	emptynorm, err := NewEmptyNormalizer()
	if err != nil {
		t.Fatalf("Failed to initialize test emptying normalization")
	}

	testRules := &rules.RuleSet{
		Normalizations: &normalization.NormalizationSequence{Seq: []normalization.Normalization{emptynorm}},
		Rules: []rules.Rule{
			{Label: "fullname", Components: []rules.RuleDemographic{
				{Label: "given", Key: "given"},
				{Label: "family", Key: "family"},
			}},
			{Label: "other", Components: []rules.RuleDemographic{
				{Label: "empty", Key: "empty"},
				{Label: "notreal", Key: "notreal"},
			}},
		},
	}

	for _, tc := range testcases {
		rawADemographics := &match.Demographics{
			Raw: tc.RawA,
		}

		hashed, e1, e2 := match.NewHashedDemographics(rawADemographics, noopHash.HashFunc, noopSalt, testRules, match.HashToString)
		if e1 != nil {
			t.Errorf("Unexpected hash error from NewHashedDemographics() %q", e1)
		}

		if e2 != nil {
			t.Fatalf("Unexpected fatal error from NewHashedDemographics() %q", e2)
		}

		result, err := match.Match(
			&match.Demographics{
				Raw: tc.RawB,
			},
			hashed,
		)
		if err != nil {
			t.Fatalf("Unexpected returned error from Match() %q", err)
		}

		// when all values normalize to an empty string it should not be possible to match
		if result.Matched {
			t.Fatalf("Expected match status to be false: %v", result)
		}
	}
}
