// Package dob trims a date of birth already in YYYY-MM-DD form and leaves
// anything else unchanged.
package dob

import (
	"fmt"
	"strings"
	"time"
)

type DOB struct{}

func New(specifier string) (*DOB, error) {
	if specifier != "dob_iso" {
		return nil, fmt.Errorf("invalid specifier: %s", specifier)
	}
	return &DOB{}, nil
}

// Apply returns the input unchanged when it cannot be read as a YYYY-MM-DD
// date. Every error on this path means the string was not a readable date;
// there is no service to fail. A record whose date of birth cannot be read
// should still hash on its remaining normalizations rather than drop every
// rule for that person (go-opm-core#1, question 3).
func (d *DOB) Apply(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if _, err := time.Parse("2006-01-02", trimmed); err == nil {
		return trimmed, nil
	}
	return input, nil
}
