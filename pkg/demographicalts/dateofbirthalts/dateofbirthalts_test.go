package dateofbirthalts_test

import (
	"slices"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
)

func TestDateOfBirthAltsFn(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"2004-02-02", []string{"2004-02-02", "2004-02-03", "2004-02-01"}},
		{"2004-02-28", []string{"2004-02-28", "2004-02-29", "2004-02-27"}},
		{"2005-02-28", []string{"2005-02-28", "2005-03-01", "2005-02-27"}},
		{"1982-08-22", []string{"1982-08-22", "1982-08-23", "1982-08-21"}},
		{"1974-05-01", []string{"1974-05-01", "1974-04-30", "1974-05-02"}},
	}

	// get the function
	fn, err := demographicalts.New("dateofbirth")
	if err != nil {
		t.Fatalf("Error creating demographic alternatives function for dateofbirth specifier: error: %s", err)
	}

	for _, tc := range cases {
		out := fn(tc.in, "test", 0)

		if len(out) != len(tc.want) {
			t.Errorf("Unexpected number of results dateofbirth demographic alternatives function: %d want: %d | %v", len(out), len(tc.want), out)
		}

		checked := make([]string, 0)
		for _, o := range out {
			if !slices.Contains(tc.want, o.Text) {
				t.Errorf("Unexpected text from dateofbirth demographic alternatives function: %s want one of %v", o.Text, tc.want)
			}

			if slices.Contains(checked, o.Text) {
				t.Errorf("Unexpected repeat text from dateofbirth demographic alternatives function: %s %v", o.Text, out)
			}
			checked = append(checked, o.Text)
		}
	}
}
