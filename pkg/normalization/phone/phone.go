// Package phone reduces a phone number to its digits, giving a ten digit
// number the US country code.
package phone

import (
	"fmt"
	"regexp"
)

var nonDigitsRegExp = regexp.MustCompile(`[^0-9]+`)

type Phone struct{}

func New(specifier string) (*Phone, error) {
	if specifier != "phone_e164" {
		return nil, fmt.Errorf("invalid specifier: %s", specifier)
	}
	return &Phone{}, nil
}

func (p *Phone) Apply(input string) (string, error) {
	digits := nonDigitsRegExp.ReplaceAllString(input, "")
	if len(digits) == 10 {
		digits = "1" + digits
	}
	return digits, nil
}
