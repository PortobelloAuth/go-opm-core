package diacriticalts_test

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
		{"Myra", []string{"Myra"}},
		{"Jones", []string{"Jones"}},
		{"San José", []string{"San José", "San Jose"}},
	}

	// get the function
	fn, err := demographicalts.New("diacritics")
	if err != nil {
		t.Fatalf("Error creating demographic alternatives function for diacritics specifier: error: %s", err)
	}

	for _, tc := range cases {
		out := fn(tc.in, "test", 0)

		if len(out) != len(tc.want) {
			t.Fatalf("Unexpected number of results diacritics demographic alternatives function: %d want: %d", len(out), len(tc.want))
		}

		checked := make([]string, 0)
		for _, o := range out {
			if !slices.Contains(tc.want, o.Text) {
				t.Fatalf("Unexpected text from diacritics demographic alternatives function: %s want one of %v", o.Text, tc.want)
			}

			if slices.Contains(checked, o.Text) {
				t.Fatalf("Unexpected repeat text from diacritics demographic alternatives function: %s %v", o.Text, out)
			}
			checked = append(checked, o.Text)
		}
	}
}
