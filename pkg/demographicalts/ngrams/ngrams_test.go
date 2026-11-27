package ngrams_test

import (
	"slices"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
)

func TestDiacriticAltsFn(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		// TODO: decide if we really want the DemographicAlts fuction to return the original string
		{"Myra", []string{"Myra", "^M", "My", "yr", "ra", "a$"}},
		{"Jones", []string{"Jones", "^J", "Jo", "on", "ne", "es", "s$"}},
		{"San José", []string{"San José", "^S", "Sa", "an", "n ", " J", "Jo", "os", "se", "e$"}},
	}

	// get the function (default to 1 deletion if the string is longer than 5 characters)
	fn, err := demographicalts.New("ngrams")
	if err != nil {
		t.Fatalf("Error creating demographic alternatives function for ngrams specifier: error: %s", err)
	}

	for _, tc := range cases {
		out := fn(tc.in, "test", 0)

		if len(out) != len(tc.want) {
			t.Errorf("Unexpected number of results ngrams demographic alternatives function: %d want: %d | %v", len(out), len(tc.want), out)
		}

		checked := make([]string, 0)
		for _, o := range out {
			if !slices.Contains(tc.want, o.Text) {
				t.Errorf("Unexpected text from ngrams demographic alternatives function: %s want one of %v", o.Text, tc.want)
			}

			if slices.Contains(checked, o.Text) {
				t.Errorf("Unexpected repeat text from ngrams demographic alternatives function: %s %v", o.Text, out)
			}
			checked = append(checked, o.Text)
		}
	}
}
