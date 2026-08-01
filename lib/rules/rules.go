package rules

import "github.com/PortobelloAuth/go-opm-core/lib/normalization"

type RuleSet struct {
	// RuleSet-level Normalizations are applied to all Demographics
	// before RuleDemographic-level Normalizations
	Normalizations normalization.NormalizationSequence
	// any Rule in a RuleSet may match
	Rules []Rule
}

type Rule struct {
	// Label describes the Rule
	Label string
	// Components of a Rule must all match
	Components []RuleDemographic
	// A Rule may contain a RuleSet or Components, not both
	Set RuleSet
}

type RuleDemographic struct {
	Normalizations normalization.NormalizationSequence
	Demographic    string
	Alternatives   DemographicAlternativesFn
}

type DemographicAlternative struct {
	// Text is the new string to use instead of the input demographic
	Text string
	// Score indicates how much the demographic was changed and therefore
	// how much matches on this alternative should be discounted. Lower
	// is better. 0 is best. Negative numbers are invalid.
	Score int
}
type DemographicAlternativesFn func(d string, limit int) []DemographicAlternative

func IdentityDemographicAlternativesFn(d string, l int) []DemographicAlternative {
	a := make([]DemographicAlternative, 1)
	a = append(a, DemographicAlternative{
		Text:  d,
		Score: 0,
	})
	return a
}
