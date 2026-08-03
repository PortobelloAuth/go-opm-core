package lowercase_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/lowercase"
)

func TestLowercaseApply(t *testing.T) {
	input := "A String With Some Uppercase Letters"
	expected := "a string with some uppercase letters"

	l, err := lowercase.New("lowercase")
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
