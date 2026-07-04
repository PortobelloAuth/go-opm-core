package punctuation_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/lib/normalization/punctuation"
)

func TestPuntuationApply(t *testing.T) {
	input := "This is a string! A #2 pencil could write it. The price is $2.00."
	expected := "This is a string A 2 pencil could write it The price is 200"

	p, err := punctuation.New("punctuation")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	output, err := p.Apply(input)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if output != expected {
		t.Errorf("output '%s' did not match expected '%s'", output, expected)
	}
}
