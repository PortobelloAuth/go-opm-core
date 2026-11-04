package phonenumbere164

import (
	"fmt"

	"github.com/nyaruka/phonenumbers"
)

type E164 struct{}

func New(specifier string) (*E164, error) {
	if specifier != "phonenumbere164" {
		return nil, fmt.Errorf("invalid specifier: %s", specifier)
	}
	return &E164{}, nil
}

func (l *E164) Apply(input string) (string, error) {
	num, err := phonenumbers.Parse(input, "US")
	if err != nil {
		return "", fmt.Errorf("failed to parse: %w", err)
	}

	// Validate against global numbering plans
	if !phonenumbers.IsValidNumber(num) {
		return "", fmt.Errorf("invalid phone number for parsed region")
	}

	// Format to strict E.164 (+15551234567)
	return phonenumbers.Format(num, phonenumbers.E164), nil
}
