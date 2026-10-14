package normalization_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization"
)

func TestAddressApply(t *testing.T) {
	cases := []struct {
		input  string
		expect string
	}{
		{
			input:  "123 Main Street Apt. 23 Norfolk, Virginia 23510",
			expect: "123 MAIN ST # 23\nNORFOLK VA 23510",
		},
		{
			input:  "a24 Calle Marginal, Bayamón, 00956, Puerto Rico",
			expect: "A24 CALLE MARGINAL\nBAYAMON PR 00956",
		},
	}

	a, err := normalization.New("address")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	for _, tc := range cases {
		output, err := a.Apply(tc.input)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if output != tc.expect {
			t.Errorf("output %q did not match expected %q", output, tc.expect)
		}
	}
}

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
