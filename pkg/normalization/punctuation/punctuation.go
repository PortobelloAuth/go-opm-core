package punctuation

import (
	"fmt"
	"regexp"
)

var punctuationRegExp = regexp.MustCompile(`[^a-zA-Z0-9\s]+`)

type Punctuation struct{}

func New(specifier string) (*Punctuation, error) {
	if specifier != "punctuation" {
		return nil, fmt.Errorf("Invalid specifier string provided: %s", specifier)
	}
	return &Punctuation{}, nil
}

func (p *Punctuation) Apply(input string) (string, error) {
	result := punctuationRegExp.ReplaceAllString(input, "")
	return result, nil
}
