package fold

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Fold struct{}

func New(specifier string) (*Fold, error) {
	if specifier != "fold" {
		return nil, fmt.Errorf("Invalid specifier string provided: %s", specifier)
	}
	return &Fold{}, nil
}

func (t *Fold) Apply(input string) (string, error) {
	return fold(input)
}

// fold implements diacritic folding by removing unicode nonspacing marks
func fold(s string) (string, error) {
	if isASCII(s) {
		return s, nil
	}
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	// None of the chained transformers can fail: invalid UTF-8 is replaced with
	// U+FFFD rather than reported. err will always be nil.
	folded, _, err := transform.String(t, s)
	return folded, err
}

func isASCII(s string) bool {
	for i := range len(s) {
		if s[i] >= utf8.RuneSelf {
			return false
		}
	}
	return true
}
