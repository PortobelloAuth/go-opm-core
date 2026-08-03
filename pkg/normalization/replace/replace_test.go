package replace_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/replace"
)

func TestReplaceApply(t *testing.T) {
	input := "A String With Some Uppercase Letters"
	expected := "A String With Six Uppercase Letters"

	r, err := replace.New("s/Some/Six/")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	output, err := r.Apply(input)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if output != expected {
		t.Errorf("output '%s' did not match expected '%s'", output, expected)
	}
}
