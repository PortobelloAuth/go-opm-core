package demographicalts_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
)

func TestIdentityDemographicAlternativesFn(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Myra", "Myra"},
		{"Jones", "Jones"},
		{"San José", "San José"},
	}

	// get the function
	fn, err := demographicalts.New("identity")
	if err != nil {
		t.Fatalf("Error creating demographic alternatives function for identity specifier: error: %s", err)
	}

	for _, tc := range cases {
		out := fn(tc.in, "test", 0)

		if len(out) != 1 {
			t.Fatalf("Identity demographic alternatives function should return 1 result: %v", out)
		}

		if out[0].Text != tc.want {
			t.Fatalf("Unexpected text from identity demographic alternatives function: %s want %s", out[0].Text, tc.want)
		}
	}
}
