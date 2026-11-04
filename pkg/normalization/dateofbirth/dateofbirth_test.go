package dateofbirth_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/dateofbirth"
)

func TestDateOfBirthApply(t *testing.T) {
	testcases := []struct {
		In   string
		Want string
	}{
		{"2007/07/17", "2007/07/17"},
		{"2007-07-17", "2007/07/17"},
		{"2007.07.17", "2007/07/17"},
		{"20070717", "2007/07/17"},
		{"17 July 2007", "2007/07/17"},
		{"17July2007", "2007/07/17"},
		{"17 Jul 2007", "2007/07/17"},
		{"17Jul2007", "2007/07/17"},
		{"July 17, 2007", "2007/07/17"},
		{"Jul 17, 2007", "2007/07/17"},
		{"2007/07", "2007/07"},
		{"2007-07", "2007/07"},
		{"2007.07", "2007/07"},
		{"July 2007", "2007/07"},
		{"Jul 2007", "2007/07"},
		{"2001", "2001"},
		{"yesterday", ""},
		{"WED", ""},
		{"next week", ""},
		{"2007/7/10", "2007/07/10"},
		{"2007-7-10", "2007/07/10"},
		{"2007.7.10", "2007/07/10"},
		{"1207", ""},
		{"2711", ""},
		{"3207.7.10", ""},
		// prefer the US formatted MM//DD/YYYY
		{"03/04/2007", "2007/03/04"},
		{"3/4/2007", "2007/03/04"},
		{"23/2/2007", "2007/02/23"},
	}

	d, err := dateofbirth.New("dateofbirth")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	for _, tc := range testcases {
		output, err := d.Apply(tc.In)
		if err != nil && tc.Want != "" {
			t.Errorf("unexpected date normalization error: %s", err)
		}

		if output != tc.Want {
			t.Errorf("output %q for %q did not match expected %q", output, tc.In, tc.Want)
		}
	}
}
