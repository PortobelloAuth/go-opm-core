package digits_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/digits"
)

func TestDigitsApply(t *testing.T) {
	input := "123-45-6789"
	expected := "123456789"

	d, err := digits.New("digits_only")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	output, err := d.Apply(input)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if output != expected {
		t.Errorf("output '%s' did not match expected '%s'", output, expected)
	}
}

func TestDigitsApplyWithWhitespace(t *testing.T) {
	input := " 6789 "
	expected := "6789"

	d, err := digits.New("digits_only")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	output, err := d.Apply(input)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if output != expected {
		t.Errorf("output '%s' did not match expected '%s'", output, expected)
	}
}
