package transliterate

import (
	"fmt"

	anyascii "github.com/anyascii/go"
)

type Transliterate struct{}

func New(specifier string) (*Transliterate, error) {
	if specifier != "transliterate" {
		return nil, fmt.Errorf("Invalid specifier string provided: %s", specifier)
	}
	return &Transliterate{}, nil
}

func (t *Transliterate) Apply(input string) (string, error) {
	result := anyascii.Transliterate(input)
	return result, nil
}
