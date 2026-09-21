package phone_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/phone"
)

func TestPhoneApply(t *testing.T) {
	input := "(212) 555-1234"
	expected := "12125551234"

	p, err := phone.New("phone_e164")
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

func TestPhoneApplyWithCountryCode(t *testing.T) {
	input := "+1 212 555 1234"
	expected := "12125551234"

	p, err := phone.New("phone_e164")
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

func TestPhoneApplyWithExtension(t *testing.T) {
	input := "212.555.1234 x12"
	expected := "212555123412"

	p, err := phone.New("phone_e164")
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

func TestPhoneApplyInternational(t *testing.T) {
	input := "+44 20 7946 0958"
	expected := "442079460958"

	p, err := phone.New("phone_e164")
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
