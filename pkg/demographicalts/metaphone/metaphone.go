package metaphone

import (
	"fmt"
	"math"

	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
	anyascii "github.com/anyascii/go"
	"github.com/vividvilla/metaphone"
)

func init() {
	demographicalts.Register("metaphone", func(s string) (demographicalts.DemographicAlternativesFn, error) {
		return Metaphone, nil
	})
	demographicalts.Register("doublemetaphone", func(s string) (demographicalts.DemographicAlternativesFn, error) {
		return DoubleMetaphone, nil
	})
}

func Metaphone(dv string, name string, l int) []demographicalts.DemographicAlternative {
	a := make([]demographicalts.DemographicAlternative, 0)
	a = append(a, demographicalts.DemographicAlternative{
		Label: name,
		Text:  dv,
		Score: 0,
	})

	transliterated := anyascii.Transliterate(dv)
	mp := metaphone.Metaphone(transliterated)

	a = append(a, demographicalts.DemographicAlternative{
		Label: fmt.Sprintf("%s|metaphone", name),
		Text:  mp,
		Score: math.MaxInt16, // NOTE: never send scores over the wire
	})

	return a
}

func DoubleMetaphone(dv string, name string, l int) []demographicalts.DemographicAlternative {
	a := make([]demographicalts.DemographicAlternative, 0)
	a = append(a, demographicalts.DemographicAlternative{
		Label: name,
		Text:  dv,
		Score: 0,
	})

	transliterated := anyascii.Transliterate(dv)
	mp1, mp2 := metaphone.DoubleMetaphone(transliterated)

	a = append(a, demographicalts.DemographicAlternative{
		Label: fmt.Sprintf("%s|doublemetaphone1", name),
		Text:  fmt.Sprintf("%s|1", mp1),
		Score: math.MaxInt16, // NOTE: never send scores over the wire
	})

	a = append(a, demographicalts.DemographicAlternative{
		Label: fmt.Sprintf("%s|doublemetaphone2", name),
		Text:  fmt.Sprintf("%s|2", mp2),
		Score: math.MaxInt16, // NOTE: never send scores over the wire
	})

	return a
}
