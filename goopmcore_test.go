package goopmcore_test

import (
	"maps"
	"testing"

	goopmcore "github.com/PortobelloAuth/go-opm-core"
	"github.com/PortobelloAuth/go-opm-core/internal/test/placeholdertestfixture"
	"github.com/PortobelloAuth/go-opm-core/pkg/match"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization"
	"github.com/PortobelloAuth/go-opm-core/pkg/rules"
)

func TestPlaintextMatch(t *testing.T) {

	check := &match.Demographics{
		Raw: map[string]string{
			"given":  "Myra",
			"family": "Jones",
		},
	}
	against_fail := &match.Demographics{
		Raw: map[string]string{
			"given":  "Adam",
			"family": "Everyman",
		},
	}

	simpleRules := &rules.RuleSet{
		Rules: []rules.Rule{
			{Label: "fullname", Components: []rules.RuleDemographic{
				{Label: "given", Key: "given"},
				{Label: "family", Key: "family"},
			}},
		},
	}

	result, err := goopmcore.PlaintextMatch(check, against_fail, simpleRules)
	if err != nil {
		t.Fatalf("Error checking mismatched demographics: %v", err)
	}

	if result {
		t.Fatalf("Match %v and %v should be false: %v", check, against_fail, result)
	}

	against_succeed := &match.Demographics{
		Raw: maps.Clone(check.Raw),
	}

	result, err = goopmcore.PlaintextMatch(check, against_succeed, simpleRules)
	if err != nil {
		t.Fatalf("Error checking matching demographics: %v", err)
	}

	if !result {
		t.Fatalf("Match %v and %v should be true: %v", check, against_fail, result)
	}
}

func TestPlaintextMatchWithPlaceholders(t *testing.T) {

	check := &match.Demographics{
		Raw: map[string]string{
			"given":     "Unknown",
			"birthDate": "01/01/1900",
			"email":     "test@test.com",
		},
	}

	// This should fail because both set's of raw hashes should
	// become empty after Placeholder removal, leaving no values
	// to match on
	against_fail := &match.Demographics{
		Raw: map[string]string{
			"given":     "Unk",
			"birthDate": "01/01/1901",
			"email":     "none@none.com",
		},
	}

	lowercase, err := normalization.New("lowercase")
	if err != nil {
		t.Fatalf("Error preparing lowercase normalizer: %v", err)
	}
	names, err := placeholdertestfixture.NewSpecified("placeholder-names")
	if err != nil {
		t.Fatalf("Error preparing names placeholder removal: %v", err)
	}
	dateofbirth, err := placeholdertestfixture.NewSpecified("placeholder-dateofbirth")
	if err != nil {
		t.Fatalf("Error preparing dateofbirth placeholder removal: %v", err)
	}
	email, err := placeholdertestfixture.NewSpecified("placeholder-email")
	if err != nil {
		t.Fatalf("Error preparing email placeholder removal: %v", err)
	}

	simpleRules := &rules.RuleSet{
		Normalizations: &normalization.NormalizationSequence{Seq: []normalization.Normalization{lowercase}},
		Rules: []rules.Rule{
			{Components: []rules.RuleDemographic{
				{Label: "given", Key: "given", Normalizations: &normalization.NormalizationSequence{Seq: []normalization.Normalization{names}}},
				{Label: "birthDate", Key: "birthDate", Normalizations: &normalization.NormalizationSequence{Seq: []normalization.Normalization{dateofbirth}}},
				{Label: "email", Key: "email", Normalizations: &normalization.NormalizationSequence{Seq: []normalization.Normalization{email}}},
			}},
		},
	}

	result, err := goopmcore.PlaintextMatch(check, against_fail, simpleRules)
	if err != nil {
		t.Fatalf("Error checking mismatched placeholder demographics: %v", err)
	}

	if result {
		t.Fatalf("Match %v and %v should be false: %v", check, against_fail, result)
	}
}
