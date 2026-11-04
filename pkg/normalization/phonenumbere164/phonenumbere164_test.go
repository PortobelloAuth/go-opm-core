package phonenumbere164_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/phonenumbere164"
)

func TestE164Apply(t *testing.T) {
	testcases := []struct {
		In   string
		Want string
	}{
		{"A String 4 you 2 test digit extraction", ""},
		{"(867)555-5309", "+18675555309"},
		{"hello", ""},
		{"12034567890", "+12034567890"},
	}

	p, err := phonenumbere164.New("phonenumbere164")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	for _, tc := range testcases {
		output, err := p.Apply(tc.In)
		if err != nil && tc.Want != "" {
			t.Errorf("unexpected phonenumbere164 normalization error: %s", err)
		}

		if output != tc.Want {
			t.Errorf("output %q for %q did not match expected %q", output, tc.In, tc.Want)
		}
	}
}

func TestE164ExtentionsApply(t *testing.T) {
	testcases := []struct {
		In   string
		Want string
	}{
		{"A String 4 you 2 test digit extraction", ""},
		{"(867)555-5309", "+18675555309"},
		{"hello", ""},
		{"12034567890", "+12034567890"},
		{"801.654.3211 ext. 55", "+18016543211 x55"},
		{"208-555-1001 x4", "+12085551001 x4"},
	}

	p, err := phonenumbere164.New("phonenumbere164-extensions")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	for _, tc := range testcases {
		output, err := p.Apply(tc.In)
		if err != nil && tc.Want != "" {
			t.Errorf("unexpected phonenumbere164 normalization error: %s", err)
		}

		if output != tc.Want {
			t.Errorf("output %q for %q did not match expected %q", output, tc.In, tc.Want)
		}
	}
}
