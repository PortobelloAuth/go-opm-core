// Package projectusat normalizes an address string per Project US@ through
// go-projectusat.
package projectusat

import (
	"fmt"

	goprojectusat "github.com/PortobelloAuth/go-projectusat"
	"github.com/PortobelloAuth/go-projectusat/pkg/address/parser"
)

// USAtNormalization applies go-projectusat's matching normalization and
// renders the result as a single line, which is the form Open Person Matching
// hashes.
//
// go-projectusat implements the standard; the parser and the reference data
// that reads a string into an address are the caller's, so that this module
// does not bind every process to a particular dataset (see the address
// README). The built-in go-projectusat parser is not implemented yet, so a
// parser is required. A verifier is optional.
type USAtNormalization struct {
	opts []goprojectusat.USAtNormalizeOption
}

func New(parse parser.ParsingFunc, verify parser.AddressVerifier) (*USAtNormalization, error) {
	if parse == nil {
		return nil, fmt.Errorf("A parser is required")
	}

	return &USAtNormalization{
		opts: []goprojectusat.USAtNormalizeOption{
			goprojectusat.WithCustomAddressParser(parse),
			goprojectusat.WithParsedAddressVerifier(verify),
			goprojectusat.WithMatchingNormalization(),
			goprojectusat.WithSingleLineFormatting(),
		},
	}, nil
}

// Apply returns the input unchanged when it cannot be read as an address.
// Every error on this path means the string was not readable or verifiable
// as a Project US@ address; there is no service to fail. A record whose
// address cannot be read should still hash on its remaining normalizations
// rather than drop every rule for that person (go-opm-core#1, question 3).
func (u *USAtNormalization) Apply(input string) (string, error) {
	result, err := goprojectusat.Normalize(input, u.opts...)
	if err != nil {
		return input, nil
	}

	return result, nil
}
