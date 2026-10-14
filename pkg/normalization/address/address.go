package address

import (
	"fmt"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/address/addresstypes"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/address/projectusat"
)

func New(specifier string) (addresstypes.AddressNormalization, error) {
	if specifier != "address" {
		return nil, fmt.Errorf("invalid specifier: %s", specifier)
	}
	pusat, err := projectusat.New(specifier)
	if err != nil {
		return nil, err
	}
	return pusat, nil
}
