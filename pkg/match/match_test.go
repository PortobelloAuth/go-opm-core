package match_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
	_ "github.com/PortobelloAuth/go-opm-core/pkg/demographicalts/diacriticalts"
	"github.com/PortobelloAuth/go-opm-core/pkg/match"
	"github.com/PortobelloAuth/go-opm-core/pkg/match/matcher/componentmatcher"
	"github.com/PortobelloAuth/go-opm-core/pkg/match/matcher/plaintextcomponent"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization"
	"github.com/PortobelloAuth/go-opm-core/pkg/rules"
)

func TestMatch(t *testing.T) {
	var testcases = []struct {
		Name   string
		Raw    map[string]string
		Hashes map[string][]string
	}{
		{
			Name: "diacritic in family name",
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
			Name: "diacritic in given and family name",
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
			Name: "no diacritics in names",
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

	plaintext, err := plaintextcomponent.New(testRules)
	if err != nil {
		t.Fatalf("Unable to instatiate plaintext Matcher: %s", err)
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			rawDemographics := &match.Demographics{
				Raw: tc.Raw,
			}

			result, err := match.Match(
				rawDemographics,
				&match.HashedDemographics{
					Hashes:  tc.Hashes,
					Matcher: plaintext,
				},
			)
			if err != nil {
				t.Fatalf("Unexpected error determining match status: %s", err)
			}

			if !result.Matched {
				t.Errorf("Expected match status to be true: %v", result)
			}

			for hkey := range tc.Hashes {
				d, ok := result.Detail[hkey]
				if !(ok && d) {
					t.Errorf("Expected match details for %s to be true: %v", hkey, result)
				}
			}
		})
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
		Name   string
		Raw    map[string]string
		Hashes map[string][]string
		Result bool
		Error  error
	}{
		{
			Name: "diacritic in family name, no error",
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
			Name: "diacritic in family name, with normailzation error",
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
			Name: "no diacritics in names, with hashing error",
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
			Name: "non-match, no diacritics in names, with hashing error",
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

	errNoop, err := NewErrorNoop()
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

	matcher, err := componentmatcher.New(
		testRules,
		errNoop.HashFunc,
		// empty salt does not matter for noop
		noopSalt,
		plaintextcomponent.PlaintextHashToString,
	)
	if err != nil {
		t.Fatalf("Unable to instatiate componentmatcher Matcher: %s", err)
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			rawDemographics := &match.Demographics{
				Raw: tc.Raw,
			}

			result, err := match.Match(
				rawDemographics,
				&match.HashedDemographics{
					Hashes:  tc.Hashes,
					Matcher: matcher,
				},
			)
			if err != nil {
				t.Fatalf("Unexpected returned error from Match() %q", err)
			}

			if !errors.Is(result.HashError, tc.Error) {
				t.Errorf("Expected match result HashError %q to be %q", err, tc.Error)
			}

			if tc.Result != result.Matched {
				t.Errorf("Expected match status to be %t: %v", tc.Result, result)
			}

			if tc.Result {
				for hkey := range tc.Hashes {
					d, ok := result.Detail[hkey]
					if !(ok && d) {
						t.Errorf("Expected match details for %s to be true: %v", hkey, result)
					}
				}
			}
		})
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

func TestMatchWithAllValuesNormalizedToBlanks(t *testing.T) {
	var testcases = []struct {
		Name string
		RawA map[string]string
		RawB map[string]string
	}{
		{
			Name: "Maria Garcia",
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
			Name: "Jose Lopez",
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
			Name: "Myra Jones",
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
			Name: "Myra Jones != Adam Smith",
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
			Name: "Empty",
			RawA: map[string]string{
				"empty": "",
			},
			RawB: map[string]string{},
		},
	}

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

	plaintext, err := plaintextcomponent.New(testRules)
	if err != nil {
		t.Fatalf("Unable to instatiate plaintextcomponent Matcher: %s", err)
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			rawADemographics := &match.Demographics{
				Raw: tc.RawA,
			}

			hashed, e1, e2 := plaintext.HashDemographics(rawADemographics)
			if e1 != nil {
				t.Errorf("Unexpected hash error from plaintext.HashDemographics() %q", e1)
			}

			if e2 != nil {
				t.Fatalf("Unexpected fatal error from plaintext.HashDemographics() %q", e2)
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
		})
	}
}

func TestNilMatcher(t *testing.T) {
	check := &match.Demographics{
		Raw: map[string]string{
			"given":  "Maria",
			"family": "García",
		},
	}

	result, err := match.Match(
		check,
		nil,
	)
	if result != nil {
		t.Errorf("Unexpected result from match.Match() against nil: %v", result)
	}
	if err == nil {
		t.Errorf("Unexpected nil error from match.Match() against nil")
	}

	lower, err := normalization.New("lowercase")
	if err != nil {
		t.Fatalf("Failed to initialize %s normalization", "lowercase")
	}

	testRules := &rules.RuleSet{
		Normalizations: &normalization.NormalizationSequence{Seq: []normalization.Normalization{lower}},
		Rules: []rules.Rule{
			{Label: "fullname", Components: []rules.RuleDemographic{
				{Label: "given", Key: "given"},
				{Label: "family", Key: "family"},
			}},
		},
	}

	plaintext, err := plaintextcomponent.New(testRules)
	if err != nil {
		t.Fatalf("Unable to instatiate plaintext Matcher: %s", err)
	}
	against, e1, e2 := plaintext.HashDemographics(check)
	if e1 != nil {
		t.Errorf("Unexpected hash error from plaintext.HashDemographics() %q", e1)
	}
	if e2 != nil {
		t.Fatalf("Unexpected fatal error from plaintext.HashDemographics() %q", e2)
	}

	// set the Matcher to nil just to test the guard
	against.Matcher = nil

	r2, err := match.Match(
		check,
		against,
	)
	if r2 != nil {
		t.Errorf("Unexpected result from match.Match() against nil: %v", r2)
	}
	if err == nil {
		t.Errorf("Unexpected nil error from match.Match() against comp with nil Matcher")
	}
}
