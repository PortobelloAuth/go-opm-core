package rules_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization"
	"github.com/PortobelloAuth/go-opm-core/pkg/rules"
)

func TestRuleSetFlatten(t *testing.T) {
	expectedLabel := "given-birthDate-email"

	lowercase, err := normalization.New("lowercase")
	if err != nil {
		t.Fatalf("Error preparing lowercase normalizer: %v", err)
	}

	simpleRules := &rules.RuleSet{
		Normalizations: &normalization.NormalizationSequence{
			Seq: []normalization.Normalization{lowercase},
		},
		Rules: []rules.Rule{
			{Components: []rules.RuleDemographic{
				{Label: "given", Key: "given"},
				{Label: "birthDate", Key: "birthDate"},
				{Label: "email", Key: "email"},
			}},
		},
	}

	flat, err := simpleRules.Flatten()
	if err != nil {
		t.Fatalf("Error flattening RuleSet: %v", err)
	}

	if len(flat) != 1 {
		t.Errorf("Expected 1 flattened rule. Got: %d. %v\n", len(flat), flat)
	}

	if flat[0].Label != expectedLabel {
		t.Errorf("Expected flattened rule label to be %q got %q\n", expectedLabel, flat[0].Label)
	}

}
