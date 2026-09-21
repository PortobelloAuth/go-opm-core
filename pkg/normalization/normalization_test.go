package normalization_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization"
)

func TestLowercaseApply(t *testing.T) {
	input := "A String With Some Uppercase Letters"
	expected := "a string with some uppercase letters"

	l, err := normalization.New("lowercase")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	output, err := l.Apply(input)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if output != expected {
		t.Errorf("output '%s' did not match expected '%s'", output, expected)
	}
}

func TestDigitsOnlyApply(t *testing.T) {
	input := "123-45-6789"
	expected := "123456789"

	d, err := normalization.New("digits_only")
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

func TestDOBIsoApply(t *testing.T) {
	input := " 1985-07-04 "
	expected := "1985-07-04"

	d, err := normalization.New("dob_iso")
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

func TestPhoneE164Apply(t *testing.T) {
	input := "(212) 555-1234"
	expected := "12125551234"

	p, err := normalization.New("phone_e164")
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

func TestUS101KeyApply(t *testing.T) {
	input := "Trần Hưng Nguyễn lives near César Chávez Avenue"
	expected := "tran hung nguyen lives near cesar chavez avenue"

	l, err := normalization.New("us-101-key")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	output, err := l.Apply(input)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if output != expected {
		t.Errorf("output '%s' did not match expected '%s'", output, expected)
	}
}

func TestUSReducedV0Apply(t *testing.T) {
	input := "Trần Hưng Nguyễn lives near César Chávez Avenue"
	expected := "tranhungnguyenlivesnearcesarchavezavenue"

	l, err := normalization.New("us-reduced_v0")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	output, err := l.Apply(input)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if output != expected {
		t.Errorf("output '%s' did not match expected '%s'", output, expected)
	}
}
