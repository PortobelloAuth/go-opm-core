package whitespace_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/whitespace"
)

func TestWhitespaceApply(t *testing.T) {
	input := "A string with\n\twhitespace"
	expected := "Astringwithwhitespace"

	w, err := whitespace.New("whitespace")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	output, err := w.Apply(input)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if output != expected {
		t.Errorf("output '%s' did not match expected '%s'", output, expected)
	}
}
