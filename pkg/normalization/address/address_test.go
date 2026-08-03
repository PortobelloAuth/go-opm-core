package address_test

import (
	"os"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/address"
)

func TestAddressNormalizationModes(t *testing.T) {
	os.Setenv("GO_OPM_CORE_ADDRESS_NORMALIZATION_MODE", "interface")

	input := "A String With Some Uppercase Letters"
	expected := "a string with some uppercase letters"

	l, err := address.New("address")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	output, err := l.Apply(input)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if output != expected {
		t.Errorf("output '%s' did not match expected '%s'", output, expected)
	}
}
