package diacriticalts

import (
	"fmt"

	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
	anyascii "github.com/anyascii/go"
)

func init() {
	demographicalts.Register("diacritics", func(s string) (demographicalts.DemographicAlternativesFn, error) {
		return DemographicAlternatives, nil
	})
}

func DemographicAlternatives(dv string, name string, l int) []demographicalts.DemographicAlternative {
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

	// sub := goprojectusat.Substitute(dv)
	// a = append(a, demographicalts.DemographicAlternative{
	// 	Label: fmt.Sprintf("%s|substitute", name),
	// 	Text:  sub,
	// 	Score: 0,
	// })
	return a
}
