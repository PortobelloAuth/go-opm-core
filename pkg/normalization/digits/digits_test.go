package digits_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/digits"
)

func TestDigitsApply(t *testing.T) {
	testcases := []struct {
		In   string
		Want string
	}{
		{"A String 4 you 2 test digit extraction", "42"},
		{"(555)867-5309", "5558675309"},
		{"hello", ""},
		{"1234567890", "1234567890"},
	}

	d, err := digits.New("digits")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	for _, tc := range testcases {
		output, err := d.Apply(tc.In)
		if err != nil {
			t.Errorf("unexpected digits normalization error: %s", err)
		}

		if output != tc.Want {
			t.Errorf("output %q for %q did not match expected %q", output, tc.In, tc.Want)
		}
	}
}
