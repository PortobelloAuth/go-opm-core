package transliterate_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/lib/normalization/transliterate"
)

func TestPuntuationApply(t *testing.T) {
	input := "άνθρωποι, René François Lacôte, Blöße, Trần Hưng Đạo, 深圳"
	expected := "anthropoi, Rene Francois Lacote, Blosse, Tran Hung Dao, ShenZhen"

	norm, err := transliterate.New("transliterate")
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
