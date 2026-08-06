package goopmcore_test

import (
	"maps"
	"testing"

	goopmcore "github.com/PortobelloAuth/go-opm-core"
	"github.com/PortobelloAuth/go-opm-core/pkg/match"
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
				{Label: "given", Demographic: "given"},
				{Label: "family", Demographic: "family"},
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
