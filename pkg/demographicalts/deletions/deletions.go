package deletions

import (
	"fmt"
	"strconv"
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

	demographicalts.Register("deletions", func(s string) (demographicalts.DemographicAlternativesFn, error) {
		parts := strings.Split(s, "-")
		if len(parts) == 3 {
			n, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("incorrect deletions specifier: number of deletions is not a number")
			}
			t, err := strconv.Atoi(parts[2])
			if err != nil {
				return nil, fmt.Errorf("incorrect deletions specifier: threshold string length is not a number")
			}

			return DeleteNAfterThreshold(n, t)
		}
		if s == "deletions" {
			return DeleteNAfterThreshold(1, 5)
		}
		return nil, fmt.Errorf("unrecogized deletions specifier")
	})
}

func DeleteNAfterThreshold(n int, threshold int) (demographicalts.DemographicAlternativesFn, error) {
	return func(dv string, name string, l int) []demographicalts.DemographicAlternative {
		a := make([]demographicalts.DemographicAlternative, 0)
		a = append(a, demographicalts.DemographicAlternative{
			Label: name,
			Text:  dv,
			Score: 0,
		})

		// don't generate alternatives for values that are too short for the threshold or the number of deletions
		if len(dv) < threshold || len(dv) <= n {
			return a
		}

		folded, err := folder.Apply(dv)
		if err == nil {
			for i := 0; i < len(folded); i++ {
				s := removeIndex(folded, i)

				if s != dv {
					a = append(a, demographicalts.DemographicAlternative{
						Label: fmt.Sprintf("%s|deletions", name),
						Text:  s,
						Score: len(dv) - len(s), // NOTE: never send scores over the wire
					})
				}
			}
		}

		return a
	}, nil
}

func removeIndex(s string, remove int) string {
	var b strings.Builder
	b.Grow(len(s) - 1)

	for i, char := range s {
		if i != remove {
			b.WriteRune(char)
		}
	}
	return b.String()
}
