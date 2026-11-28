package metaphone_test

import (
	"slices"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
)

func TestMetaphoneAltsFn(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"Myra", []string{"Myra", "MR"}},
		{"Jones", []string{"Jones", "JNS"}},
		{"San José", []string{"San José", "SN JS"}},
	}

	// get the function
	fn, err := demographicalts.New("metaphone")
	if err != nil {
		t.Fatalf("Error creating demographic alternatives function for metaphone specifier: %s", err)
	}

	for _, tc := range cases {
		out := fn(tc.in, "test", 0)

		if len(out) != len(tc.want) {
			t.Errorf("Unexpected number of results metaphone demographic alternatives function: %d want: %d | %v", len(out), len(tc.want), out)
		}

		checked := make([]string, 0)
		for _, o := range out {
			if !slices.Contains(tc.want, o.Text) {
				t.Errorf("Unexpected text from metaphone demographic alternatives function: %s want one of %v", o.Text, tc.want)
			}

			if slices.Contains(checked, o.Text) {
				t.Errorf("Unexpected repeat text from metaphone demographic alternatives function: %s %v", o.Text, out)
			}
			checked = append(checked, o.Text)
		}
	}
}

func TestDoubleMetaphoneAltsFn(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"Myra", []string{"Myra", "MR|1", "MR|2"}},
		{"Jones", []string{"Jones", "JNS|1", "ANS|2"}},
		{"San José", []string{"San José", "SNHS|1", "SNHS|2"}},
	}

	// get the function
	fn, err := demographicalts.New("doublemetaphone")
	if err != nil {
		t.Fatalf("Error creating demographic alternatives function for doublemetaphone specifier: %s", err)
	}

	for _, tc := range cases {
		out := fn(tc.in, "test", 0)

		if len(out) != len(tc.want) {
			t.Errorf("Unexpected number of results doublemetaphone demographic alternatives function: %d want: %d | %v", len(out), len(tc.want), out)
		}

		checked := make([]string, 0)
		for _, o := range out {
			if !slices.Contains(tc.want, o.Text) {
				t.Errorf("Unexpected text from doublemetaphone demographic alternatives function: %s want one of %v", o.Text, tc.want)
			}

			if slices.Contains(checked, o.Text) {
				t.Errorf("Unexpected repeat text from doublemetaphone demographic alternatives function: %s %v", o.Text, out)
			}
			checked = append(checked, o.Text)
		}
	}
}
