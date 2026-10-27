package rules

import (
	"fmt"

	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization"
)

type RuleDemographic struct {
	// TODO?: remove Label. It duplicates Key.
	Label          string
	Normalizations *normalization.NormalizationSequence
	// Key is the key by which we look up the demographic, not the value or hash itself
	Key          string
	Alternatives demographicalts.DemographicAlternativesFn
	Limit        int
}

type Rule struct {
	// Label describes the Rule
	Label string
	// Components of a Rule must all match
	Components []RuleDemographic
	// A Rule may contain a RuleSet or Components, not both
	Set *RuleSet
}

// A FlattenedRule tracks the normalizations and set of demographic alternatives that
// must match for a given permutation of a rule to match. This means there will be a
// distinct FlattenedRule for each demographic alternative of each Rule Demographic
// that exists between the root of the RuleSet and the leaf Rule Demographics.
type FlattenedRule struct {
	Label      string
	Components []RuleDemographic
}

type FlattenedRuleSet []FlattenedRule

type RuleSet struct {
	// RuleSet-level Normalizations are applied to all Demographics
	// before RuleDemographic-level Normalizations
	Normalizations *normalization.NormalizationSequence
	// any Rule in a RuleSet may match
	Rules []Rule
}

// Flatten walks the tree of Rules in the Rule Set and creates a simple
// list of Flattened Rules
func (rs *RuleSet) Flatten() (FlattenedRuleSet, error) {
	flat := make(FlattenedRuleSet, 0)

	for _, r := range rs.Rules {
		var err error
		child := make(FlattenedRuleSet, 0)
		if r.Set != nil {
			// Depth first, get the flattened rules for the next level down
			child, err = r.Set.Flatten()
			if err != nil {
				// NOTE: This is currently only here because we want to be able to return
				// errors if we find we need to.
				return nil, err
			}
		}
		layerComponents := make([]RuleDemographic, 0)
		label := r.Label
		hyphen := ""
		for _, rd := range r.Components {
			if len(r.Label) < 1 {
				label = fmt.Sprintf("%s%s", hyphen, rd.Key)
			}

			mergedseq := make([]normalization.Normalization, 0)
			if rs.Normalizations != nil && len(rs.Normalizations.Seq) > 0 {
				mergedseq = append(mergedseq, rs.Normalizations.Seq...)
			}
			if rd.Normalizations != nil && len(rd.Normalizations.Seq) > 0 {
				mergedseq = append(mergedseq, rd.Normalizations.Seq...)
			}
			flatrdnorms := &normalization.NormalizationSequence{Seq: mergedseq}

			layerComponents = append(layerComponents, RuleDemographic{
				Normalizations: flatrdnorms,
				Key:            rd.Key,
				Alternatives:   rd.Alternatives,
				Limit:          rd.Limit,
			})
			hyphen = "-"
		}
		for _, cfr := range child {
			// for every FlattenedRule in the child FlattenedRuleSet create a new FlattenedRule for
			// each component from the current RuleSet in the current FlattenedRuleSet
			flat = append(flat, FlattenedRule{
				Label:      fmt.Sprintf("%s:%s", label, cfr.Label),
				Components: append(layerComponents, cfr.Components...),
			})
		}
		if len(child) == 0 {
			// for rules with no children, just create a flattened rule
			flat = append(flat, FlattenedRule{
				Label:      fmt.Sprintf("%s", label),
				Components: layerComponents,
			})
		}
	}

	return flat, nil
}
