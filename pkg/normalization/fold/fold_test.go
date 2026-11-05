package fold_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/fold"
)

func TestFoldApply(t *testing.T) {
	input := "René François Lacôte"
	expected := "Rene Francois Lacote"

	norm, err := fold.New("fold")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	output, err := norm.Apply(input)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if output != expected {
		t.Errorf("output '%s' did not match expected '%s'", output, expected)
	}
}
