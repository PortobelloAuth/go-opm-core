package uppercase_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/uppercase"
)

func TestUppercaseApply(t *testing.T) {
	input := "A String With Some Uppercase Letters"
	expected := "A STRING WITH SOME UPPERCASE LETTERS"

	u, err := uppercase.New("uppercase")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	output, err := u.Apply(input)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if output != expected {
		t.Errorf("output '%s' did not match expected '%s'", output, expected)
	}
}
