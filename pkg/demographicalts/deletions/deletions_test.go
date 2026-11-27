package deletions_test

import (
	"slices"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
)

func TestDeletionsAltsFn(t *testing.T) {
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

func TestDeletionCount(t *testing.T) {
	cases := []struct {
		spec string
		in   string
		want int // expected length of results (excluding original)
	}{
		{"deletions-1-0", "chair", 5},  // nCk(5,1) = 5
		{"deletions-2-0", "chair", 10}, // nCk(5,2) = 10
		{"deletions-3-0", "chair", 10}, // nCk(5,2) = 10
		{"deletions-1-0", "Apple", 4},  // nCk(5,1) = 5, but the double letter means duplicates don't get included
		{"deletions-2-0", "Apple", 7},  // nCk(5,2) = 10, but the double letter means duplicates don't get included
		{"deletions-3-0", "Apple", 7},  // nCk(5,2) = 10, but the double letter means duplicates don't get included
	}

	for _, tc := range cases {
		fn, _ := demographicalts.New(tc.spec)
		out := fn(tc.in, "test", 0)

		// Subtract 1 to ignore the original string in the results
		if len(out)-1 != tc.want {
			t.Errorf("Spec %s: expected %d alternatives, got %d", tc.spec, tc.want, len(out)-1)
		}
	}
}

func TestDeletionThreshold(t *testing.T) {
	// "pear" is length 4, should produce NO alternatives.
	// "apple" is length 5, should produce alternatives.
	fn, _ := demographicalts.New("deletions-1-5")

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

func TestUnicodeDeletions(t *testing.T) {
	// "José" is 4 characters, but 'é' is 2 bytes.
	// A correct implementation should produce 4 alternatives of length 3.
	fn, _ := demographicalts.New("deletions-1-0")
	out := fn("José", "test", 0)

	for _, alt := range out {
		if alt.Text != "José" && len([]rune(alt.Text)) != 3 {
			t.Errorf("Unicode corruption: expected length 3, got %d for text %s", len([]rune(alt.Text)), alt.Text)
		}
	}
}

func TestInvalidSpecifiers(t *testing.T) {
	invalid := []string{
		"deletions-abc-5", // n is not a number
		"deletions-1-xyz", // t is not a number
	}

	for _, spec := range invalid {
		_, err := demographicalts.New(spec)
		if err == nil {
			t.Errorf("Expected error for invalid specifier %s, but got nil", spec)
		}
	}
}
