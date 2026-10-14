package address_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/address"
)

func TestApply(t *testing.T) {
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

	a, err := address.New("address")
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
