package lowercase

import (
	"fmt"
	"strings"
)

type Lowercase struct{}

func New(specifier string) (*Lowercase, error) {
	if specifier != "lowercase" {
		return nil, fmt.Errorf("invalid specifier: %s", specifier)
	}
	return &Lowercase{}, nil
}

func (l *Lowercase) Apply(input string) (string, error) {
	result := strings.ToLower(input)
	return result, nil
}
