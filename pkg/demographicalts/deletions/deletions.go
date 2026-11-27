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
	// we're using the fold normalization so that we only have 1 implementation
	// of diacritic folding in the project (not counting diacritic transliteration and
	// Pub28-specified substitution.)
	folder, err = fold.New("fold")
	if err != nil {
		panic(err)
	}

	demographicalts.Register("deletions", func(s string) (demographicalts.DemographicAlternativesFn, error) {
		parts := strings.Split(s, "-")
		n := 1
		t := 5
		if len(parts) > 1 {
			var err error
			n, err = strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("incorrect deletions specifier: number of deletions is not a number")
			}

			if n > 3 {
				return nil, fmt.Errorf("incorrect deletions specifier: more than 3 deletions is not supported")
			}

			if len(parts) > 2 {
				t, err = strconv.Atoi(parts[2])
				if err != nil {
					return nil, fmt.Errorf("incorrect deletions specifier: threshold string length is not a number")
				}
			}
		}

		return DeleteNAfterThreshold(n, t)
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

		runev := []rune(dv)

		// don't generate alternatives for values that are too short for the threshold or the number of deletions
		if len(runev) < threshold || len(runev) <= n {
			return a
		}

		folded, err := folder.Apply(dv)
		if err == nil {
			added := map[string]bool{}
			// Convert to runes to handle Unicode characters correctly
			runes := []rune(folded)

			// Generate all combinations of removing 'n' characters
			generateDeletions(runes, n, func(s string) {
				if !added[s] {
					added[s] = true
					a = append(a, demographicalts.DemographicAlternative{
						Label: fmt.Sprintf("%s|deletions", name),
						Text:  s,
						Score: len(runev) - len([]rune(s)), // NOTE: never send scores over the wire
					})
				}
			})
		}

		return a
	}, nil
}

// generateDeletions finds all combinations of removing n characters from the slice
func generateDeletions(runes []rune, n int, callback func(string)) {
	length := len(runes)
	if n <= 0 || n >= length {
		return
	}

	// indices holds the positions of the runes we will REMOVE
	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}

	for {
		// Build the string by excluding the indices in our current combination
		var b strings.Builder
		b.Grow(length - n)

		currIdx := 0
		for i := 0; i < length; i++ {
			if currIdx < n && i == indices[currIdx] {
				currIdx++
				continue
			}
			b.WriteRune(runes[i])
		}
		callback(b.String())

		// Generate the next lexicographical combination of indices
		i := n - 1
		for i >= 0 && indices[i] == length-n+i {
			i--
		}
		if i < 0 {
			break
		}
		indices[i]++
		for j := i + 1; j < n; j++ {
			indices[j] = indices[j-1] + 1
		}
	}
}
