package rules

import "github.com/PortobelloAuth/go-opm-core/lib/normalization"

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
	Alternatives   DemographicAlternativesFn
	Limit          int
}

type DemographicAlternative struct {
	// Label allows us to distinguish each demographic alternative in the resulting map
	Label string
	// Text is the new string to use instead of the input demographic
	Text string
	// Score indicates how much the demographic was changed and therefore
	// how much matches on this alternative should be discounted. Lower
	// is better. 0 is best. Negative numbers are invalid.
	Score int
}

// A DemographicAlternativesFn is a function supporting various forms of "Fuzzy" matching. A
// DemographidsAlternativeFn generates various alternative strings for the supplied demographic
// string, labeled with a derivatie of the supplied label indicating the change, and scored for
// the significance of the change.
type DemographicAlternativesFn func(demographicvalue string, label string, limit int) []DemographicAlternative

func IdentityDemographicAlternativesFn(dv string, name string, l int) []DemographicAlternative {
	a := make([]DemographicAlternative, 1)
	a = append(a, DemographicAlternative{
		Label: name,
		Text:  dv,
		Score: 0,
	})
	return a
}
