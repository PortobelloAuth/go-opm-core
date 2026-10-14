package projectusat

import (
	"fmt"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/address/addresstypes"
	goprojectusat "github.com/PortobelloAuth/go-projectusat"
	"github.com/PortobelloAuth/go-projectusat/pkg/address"
	"github.com/PortobelloAuth/go-projectusat/pkg/diacritics"
)

type ProjectUSAt struct {
	USAt *goprojectusat.USAt
}

func New(specifier string) (*ProjectUSAt, error) {
	if specifier != "address" {
		return nil, fmt.Errorf("invalid specifier: %s", specifier)
	}
	usat, err := goprojectusat.New(
		goprojectusat.WithDiacriticNormalization(diacritics.SubstituteDiacritics),
		goprojectusat.WithSecondaryAsHash(),
	)
	if err != nil {
		return nil, err
	}
	return &ProjectUSAt{
		USAt: usat,
	}, nil
}

func (p *ProjectUSAt) Apply(input string) (string, error) {
	result, err := p.USAt.Normalize(input)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (p *ProjectUSAt) ApplyAddress(input *addresstypes.Address) (*addresstypes.Address, error) {
	a := toUSAtAddress(input)
	b, err := p.USAt.NormalizeAddress(a)
	if err != nil {
		return nil, err
	}
	return fromUSAtAddress(b), nil
}

func (p *ProjectUSAt) ParseAddress(input string) (*addresstypes.Address, error) {
	result, err := p.USAt.Parse(input)
	if err != nil {
		return nil, err
	}
	return fromUSAtAddress(result), nil
}

func toUSAtAddress(src *addresstypes.Address) *address.Address {
	return &address.Address{
		BusinessName:        src.BusinessName,
		Area:                src.Area,
		PostalStation:       src.PostalStation,
		PrimaryNumber:       src.PrimaryNumber,
		Predirectional:      src.Predirectional,
		StreetName:          src.StreetName,
		StreetSuffix:        src.StreetSuffix,
		Postdirectional:     src.Postdirectional,
		SecondaryDesignator: src.SecondaryDesignator,
		SecondaryNumber:     src.SecondaryNumber,
		Detail:              src.Detail,
		City:                src.City,
		Region:              src.Region,
		Postal:              src.Postal,
		Country:             src.Country,
	}
}

func fromUSAtAddress(src *address.Address) *addresstypes.Address {
	return &addresstypes.Address{
		BusinessName:        src.BusinessName,
		Area:                src.Area,
		PostalStation:       src.PostalStation,
		PrimaryNumber:       src.PrimaryNumber,
		Predirectional:      src.Predirectional,
		StreetName:          src.StreetName,
		StreetSuffix:        src.StreetSuffix,
		Postdirectional:     src.Postdirectional,
		SecondaryDesignator: src.SecondaryDesignator,
		SecondaryNumber:     src.SecondaryNumber,
		Detail:              src.Detail,
		City:                src.City,
		Region:              src.Region,
		Postal:              src.Postal,
		Country:             src.Country,
	}
}
