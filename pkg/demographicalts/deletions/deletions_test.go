package deletions_test

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
		{"Myra", []string{"Myra"}},
		{"Jones", []string{"Jones", "ones", "Jnes", "Joes", "Jons", "Jone"}},
		{"San José", []string{"San José", "an Jose", "Sn Jose", "Sa Jose", "SanJose", "San ose", "San Jse", "San Joe", "San Jos"}},
	}

	// get the function (default to 1 deletion if the string is longer than 5 characters)
	fn, err := demographicalts.New("deletions")
	if err != nil {
		t.Fatalf("Error creating demographic alternatives function for deletions specifier: error: %s", err)
	}

	for _, tc := range cases {
		out := fn(tc.in, "test", 0)

		if len(out) != len(tc.want) {
			t.Errorf("Unexpected number of results deletions demographic alternatives function: %d want: %d | %v", len(out), len(tc.want), out)
		}

		checked := make([]string, 0)
		for _, o := range out {
			if !slices.Contains(tc.want, o.Text) {
				t.Errorf("Unexpected text from deletions demographic alternatives function: %s want one of %v", o.Text, tc.want)
			}

			if slices.Contains(checked, o.Text) {
				t.Errorf("Unexpected repeat text from deletions demographic alternatives function: %s %v", o.Text, out)
			}
			checked = append(checked, o.Text)
		}
	}
}
