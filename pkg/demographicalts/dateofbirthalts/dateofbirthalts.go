package dateofbirthalts

import (
	"fmt"
	"time"

	"github.com/PortobelloAuth/go-opm-core/pkg/demographicalts"
)

func init() {
	demographicalts.Register("dateofbirth", func(s string) (demographicalts.DemographicAlternativesFn, error) {
		return PlusOrMinus1Day, nil
	})
}

const isodate = "2006-01-02"

func PlusOrMinus1Day(dv string, name string, l int) []demographicalts.DemographicAlternative {
	a := make([]demographicalts.DemographicAlternative, 0)
	a = append(a, demographicalts.DemographicAlternative{
		Label: name,
		Text:  dv,
		Score: 0,
	})

	t, err := time.Parse(isodate, dv)
	if err != nil {
		// if we can't parse it then something else failed first; just
		// return the original value.
		// TODO: consider if this is really how we want to handle errors
		// in Demographic Alternatives functions.
		return a
	}

	plus1 := t.Add(24 * time.Hour)
	minus1 := t.Add(-24 * time.Hour)

	a = append(a, demographicalts.DemographicAlternative{
		Label: fmt.Sprintf("%s|plus1day", name),
		Text:  plus1.Format(isodate),
		Score: 0,
	})

	a = append(a, demographicalts.DemographicAlternative{
		Label: fmt.Sprintf("%s|minus1day", name),
		Text:  minus1.Format(isodate),
		Score: 0,
	})

	return a
}
