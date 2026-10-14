package addresstypes

type Address struct {
	BusinessName string // firm / business line (optional)

	// Area is the urbanization or similar neighborhood-like construct that
	// precedes the street line. USPS Publication 28 describes an urbanization
	// as "an area, sector, or development within a geographic area", which is
	// where the field gets its name.
	//
	// Project US@ requires it on Puerto Rican addresses, where it is the first
	// line of the street address block:
	//
	//	URB HIGHLAND GDNS
	//	COND LAS AMAPOLAS APT 103
	//	123 CALLE MAIN
	//
	// It is not exclusive to Puerto Rico by definition, but that is the only
	// dialect the standard requires it for today.
	Area string

	// PostalStation is the named postal station a post office box is served
	// from, such as "OLD SAN JUAN STA". USPS Publication 28 places it on the
	// line between the box and the last line, and Format writes it above the
	// street line. It is held apart from Area because it names the office
	// that handles the mail rather than a neighborhood the address is in.
	PostalStation string

	// Street line elements
	PrimaryNumber       string
	Predirectional      string
	StreetName          string
	StreetSuffix        string
	Postdirectional     string
	SecondaryDesignator string // APT, STE, ...
	SecondaryNumber     string

	// Detail is a further address number qualifying the secondary one, such as
	// the "PMB 456" of a private mailbox at a street address that already has
	// a suite. It is deliberately distinct from the secondary address parts
	// because a private mailbox must not be combined with the secondary address,
	// whereas multiple layers of secondary address units may be combined in the
	// SecondaryNumber.
	Detail string

	// Last line
	City   string
	Region string // state / province / military "state"
	Postal string // ZIP, ZIP+4, or Canadian postal code

	Country string // optional; often blank for domestic
}

type AddressNormalization interface {
	Apply(input string) (string, error)
	ApplyAddress(input *Address) (*Address, error)
	ParseAddress(input string) (*Address, error)
}
