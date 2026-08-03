package whitespace

import (
	"fmt"
	"regexp"
)

var whitespaceRegExp = regexp.MustCompile(`[\s]+`)

type Whitespace struct{}

func New(specifier string) (*Whitespace, error) {
	if specifier != "whitespace" {
		return nil, fmt.Errorf("invalid specifier: %s", specifier)
	}
	return &Whitespace{}, nil
}

func (w *Whitespace) Apply(input string) (string, error) {
	result := whitespaceRegExp.ReplaceAllString(input, "")
	return result, nil
}
