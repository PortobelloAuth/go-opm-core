package diacriticalts

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
	goprojectusat "github.com/PortobelloAuth/go-projectusat/pkg/diacritics"
	anyascii "github.com/anyascii/go"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func init() {
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

	folded := fold(dv)
	if folded != dv {
		a = append(a, demographicalts.DemographicAlternative{
			Label: fmt.Sprintf("%s|fold", name),
			Text:  folded,
			Score: 0,
		})
	}

	return a
}

// fold implements diacritic folding by removing unicode nonspacing marks
func fold(s string) string {
	if isASCII(s) {
		return s
	}
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	// None of the chained transformers can fail: invalid UTF-8 is replaced with
	// U+FFFD rather than reported.
	folded, _, _ := transform.String(t, s)
	return folded
}

func isASCII(s string) bool {
	for i := range len(s) {
		if s[i] >= utf8.RuneSelf {
			return false
		}
	}
	return true
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
