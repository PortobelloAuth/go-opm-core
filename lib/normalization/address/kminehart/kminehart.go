package kminehart

import (
	"fmt"

	boostport "github.com/Boostport/address"
	parse_address "github.com/kminehart/golang-address"
)

type KMinehart struct{}

func New(specifier string) (*KMinehart, error) {
	if specifier != "address" {
		return nil, fmt.Errorf("invalid specifier: %s", specifier)
	}
	return &KMinehart{}, nil
}

func (l *KMinehart) Apply(input string) (string, error) {
	parsed, err := parse_address.Normalize(input)
	if err != nil {
		return "", fmt.Errorf("Unable to parse address: %w", err)
	}

	/*
		type Address struct {
			House           int
			StreetDirection string
			StreetName      string
			StreetType      string
			SuiteType       string
			SuiteNumber     string
		}
	*/

	validated, err := boostport.NewValid(
		boostport.WithStreetAddress([]string{
			parsed.StreetName,
		}),
		boostport.WithLocality("Melbourne"),
		boostport.WithAdministrativeArea("VIC"),
		boostport.WithPostCode("3000"),
	)
	return result, nil
}
