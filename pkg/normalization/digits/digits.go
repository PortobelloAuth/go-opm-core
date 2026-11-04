package digits

import (
	"fmt"
	"strings"

	"golang.org/x/text/runes"
)

var notAsciiDigit = runes.Predicate(func(r rune) bool {
	return r < '0' || r > '9'
})

type Digits struct{}

func New(specifier string) (*Digits, error) {
	if specifier != "digits" {
		return nil, fmt.Errorf("invalid specifier: %s", specifier)
	}
	return &Digits{}, nil
}

func (l *Digits) Apply(input string) (string, error) {
	var b strings.Builder
	b.Grow(len(input))
	for _, r := range input {
		if notAsciiDigit.Contains(r) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String(), nil
}
