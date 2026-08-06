package demographicalts

import (
	"fmt"
	"strings"
)

type DemographicAlternativesConstructor func(specifier string) (DemographicAlternativesFn, error)

var altFnRegistery map[string]DemographicAlternativesConstructor
var newCalled = false

func init() {
	// NOTE: NetworkSettings implementations should register themselves
	altFnRegistery = make(map[string]DemographicAlternativesConstructor)
	Register("identity", func(specifier string) (DemographicAlternativesFn, error) {
		return IdentityDemographicAlternativesFn, nil
	})
}

func Register(specifier string, fn DemographicAlternativesConstructor) error {
	// FIXME?: don't allow this to be called after New()? It makes testing annoying.
	if newCalled {
		return fmt.Errorf("Register() called after New()")
	}

	// FIXME?: don't allow overwriting builtins? If not, how would a package user test out replacement code?
	spec := strings.ToLower(specifier)
	// _, found := netSettingsRegistery[spec]
	// if found {
	// 	return fmt.Errorf("Registration already exists for %s", spec)
	// }
	altFnRegistery[spec] = fn

	return nil
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
// string, labeled with a derivative of the supplied label indicating the change, and scored for
// the significance of the change. The `limit` indicates that only alternatives with a score
// lower than the limit should be returned.
type DemographicAlternativesFn func(demographicvalue string, label string, limit int) []DemographicAlternative

// TODO: determine if we should migrate to an interface like our other registries
// func (dafn DemographicAlternativesFn) Apply(d string, lbl string, lmt int) []DemographicAlternative {
// 	return dafn(d, lbl, lmt)
// }

// type DemographicAltsFn interface {
// 	Apply() func(demographicvalue string, label string, limit int) []DemographicAlternative
// }

func IdentityDemographicAlternativesFn(dv string, name string, l int) []DemographicAlternative {
	a := make([]DemographicAlternative, 0)
	a = append(a, DemographicAlternative{
		Label: name,
		Text:  dv,
		Score: 0,
	})
	return a
}

func New(specifier string) (DemographicAlternativesFn, error) {
	newCalled = true

	sep := "-"
	parts := strings.Split(strings.ToLower(specifier), sep)

	for len(parts) > 0 {
		prefix := strings.Join(parts, sep)
		fn, ok := altFnRegistery[prefix]
		if ok {
			return fn(specifier)
		}
		parts = parts[0 : len(parts)-1]
	}

	return nil, fmt.Errorf("%s is not a supported demographic alternatives function type", specifier)
}
