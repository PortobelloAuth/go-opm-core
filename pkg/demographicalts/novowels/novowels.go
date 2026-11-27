package novowels

import (
	"fmt"
	"strings"

	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/fold"
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

	demographicalts.Register("novowels", func(s string) (demographicalts.DemographicAlternativesFn, error) {
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

	folded, err := folder.Apply(dv)
	if err == nil {
		s := removeVowels(folded)

		if s != dv {
			a = append(a, demographicalts.DemographicAlternative{
				Label: fmt.Sprintf("%s|novowells", name),
				Text:  s,
				Score: len(dv) - len(s), // NOTE: never send scores over the wire
			})
		}
	}

	return a
}

func removeVowels(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	for _, char := range s {
		switch char {
		case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
			continue
		default:
			b.WriteRune(char)
		}
	}
	return b.String()
}
