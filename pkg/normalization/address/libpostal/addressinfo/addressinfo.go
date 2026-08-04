package addressinfo

import (
	"regexp"
	"strings"

	anyascii "github.com/anyascii/go"
)

var whitespaceRegExp = regexp.MustCompile(`[\s]+`)

type Address struct {
	BuisinessName string

	// Street line elements
	PrimaryNumber       string
	Predirectional      string
	StreetName          string
	StreetSuffix        string
	Postdirectional     string
	SecondaryDesignator string // APT, STE, ...
	SecondaryNumber     string

	// Last line
	City   string
	Region string // state / province / military "state"
	Postal string // ZIP, ZIP+4, or Canadian postal code

	Country string // optional; often blank for domestic
}

func (a *Address) Format() string {
	street := strings.Join([]string{
		a.PrimaryNumber,
		a.Predirectional,
		a.StreetName,
		a.StreetSuffix,
		a.Postdirectional,
		a.SecondaryDesignator,
		a.SecondaryNumber,
	}, " ")
	street = whitespaceRegExp.ReplaceAllString(street, " ")

	cityline := strings.Join([]string{
		a.City,
		a.Region,
		a.Postal,
	}, " ")

	multiline := []string{
		street,
		cityline,
	}

	if len(a.BuisinessName) > 0 {
		multiline = append([]string{a.BuisinessName}, multiline...)
	}

	if len(a.Country) > 0 {
		multiline = append(multiline, a.Country)
	}

	joined := strings.Join(multiline, "\n")

	// TODO: this should use go-projectusat's diacritics.Substitute
	transliterated := anyascii.Transliterate(joined)

	return strings.ToUpper(transliterated)
}
