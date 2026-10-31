package placeholder_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/placeholder"
)

func TestBadSpecifierPlaceholder(t *testing.T) {
	cases := []struct {
		Label string
		Spec  string
	}{
		{"No specifier", ""},
		{"Incorrect specifier", "specifier"},
		{"Unknown list", "placeholder-not a real list"},
		{"Incorrect list seporator", "placeholder/names"},
	}

	for _, tc := range cases {
		t.Run(tc.Label, func(t *testing.T) {
			_, err := placeholder.New(tc.Spec)
			if err == nil {
				t.Errorf("Expected error instantiating %q placeholder normalizer", tc.Spec)
			}
		})
	}

}

func TestNamesPlaceholder(t *testing.T) {
	nameCases := []struct {
		label string
		In    string
		Want  string
		Ctx   map[string]string
	}{
		{"Baby Boy, no family name", "Baby Boy", "Baby Boy", map[string]string{"family": ""}},
		{"Baby Boy, with family name", "Baby Boy", "", map[string]string{"family": "Doe"}},
		{"Baby Girl, no family name", "Baby Girl", "Baby Girl", map[string]string{}},
		{"Baby Girl, with family name", "Baby Girl", "", map[string]string{"family": "Doe"}},
		{"Twin A", "Twin A", "", map[string]string{}},
		{"Twin B", "Twin B", "", map[string]string{"family": "Doe"}},
		{"Newborn", "Newborn", "", map[string]string{}},
		{"Not Applicable", "N/A", "", map[string]string{}},
		{"None", "None", "", map[string]string{}},
		{"To be determined", "TBD", "", map[string]string{}},
		{"John Doe", "John", "", map[string]string{"family": "Doe"}},
		{"Jane Doe", "JANE", "", map[string]string{"family": "Doe"}},
		{"John with other last name", "John", "John", map[string]string{"family": "Smith"}},
		{"Jane with other last name", "JANE", "JANE", map[string]string{"family": "Jones"}},
		{"Other first name with last name Doe", "Mary", "Mary", map[string]string{"family": "Doe"}},
		{"A first name", "Hector", "Hector", map[string]string{}},
	}

	p, err := placeholder.New("placeholder-names")
	if err != nil {
		t.Fatalf("Unable to instantiate 'names' placeholder normalizer")
	}
	for _, tc := range nameCases {
		t.Run(tc.label, func(t *testing.T) {
			n := p.WithContext(func(s string) (string, bool) {
				v, ok := tc.Ctx[s]
				return v, ok
			})
			got, err := n.Apply(tc.In)
			if err != nil {
				t.Errorf("Error using placeholder normalizer for %q", tc.In)
			}
			if got != tc.Want {
				t.Errorf("'names' placeholder normalizer returned %q; want %q", got, tc.Want)
			}
		})
	}
}

func TestDateOfBirthPlaceholder(t *testing.T) {
	testCases := []struct {
		label string
		In    string
		Want  string
		Ctx   map[string]string
	}{
		{"1 Jan 1900", "01/01/1900", "", map[string]string{}},
		{"1 Jan 1901", "01/01/1901", "", map[string]string{}},
		{"1 Jan 1800", "01/01/1800", "", map[string]string{}},
		{"1 Jan 2000", "01/01/2000", "", map[string]string{}},
		{"11 Nov 1111", "11/11/1111", "", map[string]string{}},
		{"1 Jan 0001", "01/01/0001", "", map[string]string{}},
		{"1 Feb 1987", "01/02/1987", "01/02/1987", map[string]string{}},
	}

	p, err := placeholder.New("placeholder-dateofbirth")
	if err != nil {
		t.Fatalf("Unable to instantiate 'dateofbirth' placeholder normalizer")
	}
	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			n := p.WithContext(func(s string) (string, bool) {
				v, ok := tc.Ctx[s]
				return v, ok
			})
			got, err := n.Apply(tc.In)
			if err != nil {
				t.Errorf("Error using placeholder normalizer for %q", tc.In)
			}
			if got != tc.Want {
				t.Errorf("'dateofbirth' placeholder normalizer returned %q; want %q", got, tc.Want)
			}
		})
	}
}

func TestLast4Placeholder(t *testing.T) {
	testCases := []struct {
		label string
		In    string
		Want  string
		Ctx   map[string]string
	}{
		{"zeros", "0000", "", map[string]string{}},
		{"ones", "1111", "", map[string]string{}},
		{"twos", "2222", "", map[string]string{}},
		{"threes", "3333", "", map[string]string{}},
		{"fours", "4444", "", map[string]string{}},
		{"fives", "5555", "", map[string]string{}},
		{"sixes", "6666", "", map[string]string{}},
		{"sevens", "7777", "", map[string]string{}},
		{"eights", "8888", "", map[string]string{}},
		{"nines", "9999", "", map[string]string{}},
		{"1234", "1234", "", map[string]string{}},
		{"non-repeating", "5198", "5198", map[string]string{}},
		{"non-digit", "5a98", "", map[string]string{}},
		{"short", "598", "", map[string]string{}},
		{"long", "59872", "", map[string]string{}},
	}

	p, err := placeholder.New("placeholder-last4")
	if err != nil {
		t.Fatalf("Unable to instantiate 'last4' placeholder normalizer")
	}
	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			n := p.WithContext(func(s string) (string, bool) {
				v, ok := tc.Ctx[s]
				return v, ok
			})
			got, err := n.Apply(tc.In)
			if err != nil {
				t.Errorf("Error using placeholder normalizer for %q", tc.In)
			}
			if got != tc.Want {
				t.Errorf("'last4' placeholder normalizer returned %q; want %q", got, tc.Want)
			}
		})
	}
}

func TestPhonePlaceholder(t *testing.T) {
	testCases := []struct {
		label string
		In    string
		Want  string
		Ctx   map[string]string
	}{
		{"formatted zeros", "000-000-0000", "", map[string]string{}},
		// {"unformatted zeros", "0000000000", "", map[string]string{}},
		// {"parentheses formatted zeros", "(000)000-0000", "", map[string]string{}},
		// {"dot-formatted zeros", "000.000.0000", "", map[string]string{}},
		{"ones", "111-111-1111", "", map[string]string{}},
		{"fives", "555-555-5555", "", map[string]string{}},
		{"1234", "123-456-7890", "", map[string]string{}},
		{"555 exchange", "555-218-2342", "", map[string]string{}},
		// {"short", "218-2342", "", map[string]string{}},
		// {"long", "123-887-218-2342", "", map[string]string{}},
		// NOTE: The placeholder standard is silent on short numbers, long numbers,
		// extensions, and international numbers. e164 digits...?
		// {"US 10-digit", "480-218-2342", "480-218-2342", map[string]string{}},
		// {"US 11-digit", "1-480-218-2342", "1-480-218-2342", map[string]string{}},
	}

	p, err := placeholder.New("placeholder-phonenumber")
	if err != nil {
		t.Fatalf("Unable to instantiate 'phonenumber' placeholder normalizer")
	}
	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			n := p.WithContext(func(s string) (string, bool) {
				v, ok := tc.Ctx[s]
				return v, ok
			})
			got, err := n.Apply(tc.In)
			if err != nil {
				t.Errorf("Error using placeholder normalizer for %q", tc.In)
			}
			if got != tc.Want {
				t.Errorf("'phonenumber' placeholder normalizer returned %q; want %q", got, tc.Want)
			}
		})
	}
}

func TestEMailPlaceholder(t *testing.T) {
	testCases := []struct {
		label string
		In    string
		Want  string
		Ctx   map[string]string
	}{
		{"test at test.com", "test@test.com", "", map[string]string{}},
		{"unknown at unknown.com", "unknown@unknown.com", "", map[string]string{}},
		{"none at none.com", "none@none.com", "", map[string]string{}},
		{"noreply", "noreply@portobello.health", "", map[string]string{}},
		// {"hyphenated noreply", "no-reply@portobello.health", "", map[string]string{}},
		// {"dotted noreply", "no.reply@portobello.health", "", map[string]string{}},
		// {"underscored noreply", "no_reply@portobello.health", "", map[string]string{}},
		{"example.com", "fake@example.com", "", map[string]string{}},
		{"example.com", "abe@example.com", "", map[string]string{}},
	}

	p, err := placeholder.New("placeholder-email")
	if err != nil {
		t.Fatalf("Unable to instantiate 'email' placeholder normalizer")
	}
	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			n := p.WithContext(func(s string) (string, bool) {
				v, ok := tc.Ctx[s]
				return v, ok
			})
			got, err := n.Apply(tc.In)
			if err != nil {
				t.Errorf("Error using placeholder normalizer for %q", tc.In)
			}
			if got != tc.Want {
				t.Errorf("'email' placeholder normalizer returned %q; want %q", got, tc.Want)
			}
		})
	}
}

func TestAddressPlaceholder(t *testing.T) {
	testCases := []struct {
		label string
		In    string
		Want  string
		Ctx   map[string]string
	}{
		{"123 Main St", "123 Main St", "", map[string]string{}},
		{"123 MAIN ST", "123 MAIN ST", "", map[string]string{}},
		{"PO Box 0", "PO Box 0", "", map[string]string{}},
		{"Unknown", "Unknown", "", map[string]string{}},
		{"UNKNOWN", "UNKNOWN", "", map[string]string{}},
		{"Homeless", "Homeless", "", map[string]string{}},
		{"No Fixed Address", "No Fixed Address", "", map[string]string{}},
	}

	p, err := placeholder.New("placeholder-address")
	if err != nil {
		t.Fatalf("Unable to instantiate 'address' placeholder normalizer")
	}
	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			n := p.WithContext(func(s string) (string, bool) {
				v, ok := tc.Ctx[s]
				return v, ok
			})
			got, err := n.Apply(tc.In)
			if err != nil {
				t.Errorf("Error using placeholder normalizer for %q", tc.In)
			}
			if got != tc.Want {
				t.Errorf("'address' placeholder normalizer returned %q; want %q", got, tc.Want)
			}
		})
	}
}

func TestSubscriberIDPlaceholder(t *testing.T) {
	testCases := []struct {
		label string
		In    string
		Want  string
		Ctx   map[string]string
	}{
		{"Pending", "Pending", "", map[string]string{}},
		{"PENDING", "PENDING", "", map[string]string{}},
		{"TBD", "TBD", "", map[string]string{}},
		{"None", "None", "", map[string]string{}},
		{"none", "none", "", map[string]string{}},
		{"zeros", "00000000000", "", map[string]string{}},
		{"ones", "111111111", "", map[string]string{}},
		{"nines", "999999999999999", "", map[string]string{}},
	}

	p, err := placeholder.New("placeholder-subscriberid")
	if err != nil {
		t.Fatalf("Unable to instantiate 'subscriberid' placeholder normalizer")
	}
	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			n := p.WithContext(func(s string) (string, bool) {
				v, ok := tc.Ctx[s]
				return v, ok
			})
			got, err := n.Apply(tc.In)
			if err != nil {
				t.Errorf("Error using placeholder normalizer for %q", tc.In)
			}
			if got != tc.Want {
				t.Errorf("'subscriberid' placeholder normalizer returned %q; want %q", got, tc.Want)
			}
		})
	}
}
