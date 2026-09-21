// Package digits keeps only the digit characters 0-9 in the input.
package digits

import (
	"fmt"
	"regexp"
)

var digitsRegExp = regexp.MustCompile(`[^0-9]+`)

type Digits struct{}

func New(specifier string) (*Digits, error) {
	if specifier != "digits_only" {
		return nil, fmt.Errorf("invalid specifier: %s", specifier)
	}
	return &Digits{}, nil
}

func (d *Digits) Apply(input string) (string, error) {
	result := digitsRegExp.ReplaceAllString(input, "")
	return result, nil
}
