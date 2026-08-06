package rules

import (
	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization"
)

type RuleSet struct {
	// RuleSet-level Normalizations are applied to all Demographics
	// before RuleDemographic-level Normalizations
	Normalizations *normalization.NormalizationSequence
	// any Rule in a RuleSet may match
	Rules []Rule
}

type Rule struct {
	// Label describes the Rule
	Label string
	// Components of a Rule must all match
	Components []RuleDemographic
	// A Rule may contain a RuleSet or Components, not both
	Set *RuleSet
}

type RuleDemographic struct {
	Label          string
	Normalizations *normalization.NormalizationSequence
	Demographic    string
	Alternatives   demographicalts.DemographicAlternativesFn
	Limit          int
}
