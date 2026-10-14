package projectusat_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/address/addresstypes"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/address/projectusat"
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

	a, err := projectusat.New("address")
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

func TestParse(t *testing.T) {
	cases := []struct {
		input  string
		expect *addresstypes.Address
	}{
		{
			input: "123 Main Street Apt. 23 Norfolk, Virginia 23510",
			expect: &addresstypes.Address{
				BusinessName:    "",
				Area:            "",
				PostalStation:   "",
				PrimaryNumber:   "123",
				Predirectional:  "",
				StreetName:      "MAIN",
				StreetSuffix:    "ST",
				Postdirectional: "",
				// NOTE: Parse() does not apply normalization options
				// like WithSecondaryAsHash()
				SecondaryDesignator: "APT",
				SecondaryNumber:     "23",
				Detail:              "",
				City:                "NORFOLK",
				Region:              "VA",
				Postal:              "23510",
				Country:             "",
			},
		},
		{
			input: "a24 Calle Marginal, Bayamón, 00956, Puerto Rico",
			expect: &addresstypes.Address{
				BusinessName:        "",
				Area:                "",
				PostalStation:       "",
				PrimaryNumber:       "A24",
				Predirectional:      "",
				StreetName:          "CALLE MARGINAL",
				StreetSuffix:        "",
				Postdirectional:     "",
				SecondaryDesignator: "",
				SecondaryNumber:     "",
				Detail:              "",
				// NOTE: Parse() does not apply normalization options
				// like WithDiacriticNormalization()
				City:    "BAYAMÓN",
				Region:  "PR",
				Postal:  "00956",
				Country: "",
			},
		},
	}

	a, err := projectusat.New("address")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	for _, tc := range cases {
		output, err := a.ParseAddress(tc.input)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if *output != *tc.expect {
			t.Errorf("output %q did not match expected %q", output, tc.expect)
		}
	}
}
