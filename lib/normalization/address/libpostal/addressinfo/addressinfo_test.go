package addressinfo_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/lib/normalization/address/libpostal/addressinfo"
)

func TestFormat(t *testing.T) {
	addr := &addressinfo.Address{
		PrimaryNumber:       "123",
		StreetName:          "Main",
		StreetSuffix:        "Street",
		SecondaryDesignator: "Apartment",
		SecondaryNumber:     "4",
		City:                "Springfield",
		Region:              "Illinois",
		Postal:              "62701",
	}
	want := "123 MAIN ST APT 4\nSPRINGFIELD IL 62701"

	formatted := addr.Format()
	if formatted != want {
		t.Fatalf("Format() return unexpected address string: %s\nwant: %s", formatted, want)
	}
}
