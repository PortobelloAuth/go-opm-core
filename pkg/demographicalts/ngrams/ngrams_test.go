package ngrams_test

import (
	"slices"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
)

func TestNGramsAltsFn(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		// TODO: decide if we really want the DemographicAlts fuction to return the original string
		{"Myra", []string{"Myra", "^M", "My", "yr", "ra", "a$"}},
		{"Jones", []string{"Jones", "^J", "Jo", "on", "ne", "es", "s$"}},
		{"San José", []string{"San José", "^S", "Sa", "an", "n ", " J", "Jo", "os", "se", "e$"}},
	}

	// get the function (default to bi-grams)
	fn, err := demographicalts.New("ngrams")
	if err != nil {
		t.Fatalf("Error creating demographic alternatives function for ngrams specifier: %s", err)
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

func TestNGrams32AltsFn(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"Amy", []string{"Amy", "^^A", "^Am", "Amy", "my$", "y$$"}},
		{"Jones", []string{"Jones", "^^J", "^Jo", "Jon", "one", "nes", "es$", "s$$"}},
		{"San José", []string{"San José", "^^S", "^Sa", "San", "an ", "n J", " Jo", "Jos", "ose", "se$", "e$$"}},
	}

	// get the function (default to 1 deletion if the string is longer than 5 characters)
	fn, err := demographicalts.New("ngrams-3-2")
	if err != nil {
		t.Fatalf("Error creating demographic alternatives function for ngrams specifier: %q error: %s", "ngrams-3-2", err)
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

			if o.Score != 0 {
				// don't include the original demographic value in our repeat check because for ngrams we may want an
				// ngram that covers the whole original value.
				checked = append(checked, o.Text)
			}
		}
	}
}

func TestNGramThreshold(t *testing.T) {
	// "pear" is length 4, should produce NO alternatives.
	// "apple" is length 5, should produce alternatives.
	fn, err := demographicalts.New("ngrams-2-5")
	if err != nil {
		t.Fatalf("Error creating demographic alternatives function for ngrams specifier: %s error: %s", "ngrams-2-5", err)
	}

	// Case: Under threshold
	outShort := fn("pear", "test", 0)
	if len(outShort) != 1 {
		t.Errorf("Expected only original string for input below threshold, got %d", len(outShort))
	}

	// Case: At threshold
	outAt := fn("apple", "test", 0)
	if len(outAt) <= 1 {
		t.Errorf("Expected alternatives for input at threshold, got only %d", len(outAt))
	}
}

func TestInvalidNGramSpecifiers(t *testing.T) {
	invalid := []string{
		"ngram-abc-5", // n is not a number
		"ngram-2-xyz", // t is not a number
		"ngram-1-5",   // n is less than 2
	}

	for _, spec := range invalid {
		_, err := demographicalts.New(spec)
		if err == nil {
			t.Errorf("Expected error for invalid specifier %s, but got nil", spec)
		}
	}
}
