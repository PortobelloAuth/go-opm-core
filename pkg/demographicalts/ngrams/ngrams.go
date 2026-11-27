package ngrams

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/fold"
)

var folder *fold.Fold

func init() {
	var err error
	// we're using the fold normalization so that we only have 1 implementation
	// of diacritic folding in the project (not counting diacritic transliteration and
	// Pub28-specified substitution.)
	folder, err = fold.New("fold")
	if err != nil {
		panic(err)
	}

	demographicalts.Register("ngrams", func(s string) (demographicalts.DemographicAlternativesFn, error) {
		parts := strings.Split(s, "-")
		if len(parts) >= 2 {
			t := 0
			n, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("incorrect ngrams specifier: length of ngrams is not a number")
			}
			if n < 2 {
				return nil, fmt.Errorf("incorrect ngrams specifier: length of ngrams is less than 2")
			}

			if len(parts) > 2 {
				tspec, err := strconv.Atoi(parts[2])
				if err != nil {
					return nil, fmt.Errorf("incorrect ngrams specifier: threshold string length is not a number")
				}
				t = tspec
			}

			return NGrams(n, t)
		}
		if s == "ngrams" {
			return NGrams(2, 0)
		}
		return nil, fmt.Errorf("unrecognized ngrams specifier")
	})
}

func NGrams(n int, threshold int) (demographicalts.DemographicAlternativesFn, error) {
	return func(dv string, name string, l int) []demographicalts.DemographicAlternative {
		a := make([]demographicalts.DemographicAlternative, 0)
		a = append(a, demographicalts.DemographicAlternative{
			Label: name,
			Text:  dv,
			Score: 0,
		})

		runev := []rune(dv)

		// don't generate alternatives for values that are too short for the threshold
		if len(runev) < threshold {
			return a
		}

		folded, err := folder.Apply(dv)
		if err == nil {
			added := map[string]bool{}
			runef := []rune(folded)
			for i := (1 - n); i < len(runef); i++ {
				s := ngram(runef, i, n)

				_, found := added[s]
				if !found {
					added[s] = true
					a = append(a, demographicalts.DemographicAlternative{
						Label: fmt.Sprintf("%s|%d-grams", name, n),
						Text:  s,
						Score: math.MaxInt16, // NOTE: never send scores over the wire
					})
				}
			}
		}

		return a
	}, nil
}

func ngram(s []rune, start, length int) string {
	var b strings.Builder
	b.Grow(length)
	runeswritten := 0

	if start < (1-length) || start >= len(s) {
		return ""
	}
	end := start + length

	// add placeholders for before the start
	for s := start; s < 0; s++ {
		b.WriteRune('^')
		runeswritten++
	}

	for i, char := range s {
		if i >= start && i < end {
			b.WriteRune(char)
			runeswritten++
		}
		if runeswritten >= length {
			break
		}
	}

	// add placeholders for after the end
	for runeswritten < length {
		b.WriteRune('$')
		runeswritten++
	}

	return b.String()
}
