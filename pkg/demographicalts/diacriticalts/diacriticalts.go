package diacriticalts

import (
	"fmt"

	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/fold"
	goprojectusat "github.com/PortobelloAuth/go-projectusat/pkg/diacritics"
	anyascii "github.com/anyascii/go"
)

var folder *fold.Fold

func init() {
	var err error
	// we're using the fold normalization so that we only have 1 implmentation
	// of diacritic folding in the project (not counting diacritic tranliteration and
	// Pub28-specified substitution.)
	folder, err = fold.New("fold")
	if err != nil {
		panic(err)
	}
	demographicalts.Register("diacritics", func(s string) (demographicalts.DemographicAlternativesFn, error) {
		return DemographicAlternatives, nil
	})
	demographicalts.Register("diacritics-fold", func(s string) (demographicalts.DemographicAlternativesFn, error) {
		return Fold, nil
	})
	demographicalts.Register("diacritics-transliterate", func(s string) (demographicalts.DemographicAlternativesFn, error) {
		return Transliterate, nil
	})
	demographicalts.Register("diacritics-pub28substitute", func(s string) (demographicalts.DemographicAlternativesFn, error) {
		return Pub28Substitute, nil
	})
}

func DemographicAlternatives(dv string, name string, l int) []demographicalts.DemographicAlternative {
	a := make([]demographicalts.DemographicAlternative, 0)
	a = append(a, Fold(dv, name, l)...)
	a = append(a, Transliterate(dv, name, l)...)
	// we don't include Pub28Substitute() because it also lowercases the value, which frequently causes it
	// to change in ways this package isn't actually interested in.
	return deduplicate(a)
}

func Fold(dv string, name string, l int) []demographicalts.DemographicAlternative {
	a := make([]demographicalts.DemographicAlternative, 0)
	a = append(a, demographicalts.DemographicAlternative{
		Label: name,
		Text:  dv,
		Score: 0,
	})

	folded, err := folder.Apply(dv)
	// err should never NOT be nil, but we'll be careful here anyway
	if err == nil && folded != dv {
		a = append(a, demographicalts.DemographicAlternative{
			Label: fmt.Sprintf("%s|fold", name),
			Text:  folded,
			Score: 0,
		})
	}

	return a
}

func Transliterate(dv string, name string, l int) []demographicalts.DemographicAlternative {
	a := make([]demographicalts.DemographicAlternative, 0)
	a = append(a, demographicalts.DemographicAlternative{
		Label: name,
		Text:  dv,
		Score: 0,
	})

	transliterated := anyascii.Transliterate(dv)
	if transliterated != dv {
		a = append(a, demographicalts.DemographicAlternative{
			Label: fmt.Sprintf("%s|transliterate", name),
			Text:  transliterated,
			Score: 0,
		})
	}

	return a
}

func Pub28Substitute(dv string, name string, l int) []demographicalts.DemographicAlternative {
	a := make([]demographicalts.DemographicAlternative, 0)
	a = append(a, demographicalts.DemographicAlternative{
		Label: name,
		Text:  dv,
		Score: 0,
	})

	sub, err := goprojectusat.Substitute(dv)
	if err == nil && sub != dv {
		a = append(a, demographicalts.DemographicAlternative{
			Label: fmt.Sprintf("%s|pub28substitute", name),
			Text:  sub,
			Score: 0,
		})
	}

	return a
}

func deduplicate(in []demographicalts.DemographicAlternative) []demographicalts.DemographicAlternative {
	have := map[string]bool{}
	a := make([]demographicalts.DemographicAlternative, 0)
	for _, alt := range in {
		if !have[alt.Text] {
			a = append(a, alt)
			have[alt.Text] = true
		}
	}

	return a
}
