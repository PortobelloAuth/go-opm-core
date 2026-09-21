package dob_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/dob"
)

func TestDOBApply(t *testing.T) {
	input := " 1985-07-04 "
	expected := "1985-07-04"

	d, err := dob.New("dob_iso")
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

func TestDOBApplyWithUnreadableFormat(t *testing.T) {
	input := "07/04/1985"
	expected := "07/04/1985"

	d, err := dob.New("dob_iso")
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

func TestDOBApplyWithInvalidDate(t *testing.T) {
	input := "1985-13-40"
	expected := "1985-13-40"

	d, err := dob.New("dob_iso")
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

func TestDOBApplyWithEmptyInput(t *testing.T) {
	input := ""
	expected := ""

	d, err := dob.New("dob_iso")
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
