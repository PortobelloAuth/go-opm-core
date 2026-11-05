package dateofbirth

import (
	"fmt"
	"time"
)

type DateOfBirth struct{}

var supportedDateLayouts = map[string][]string{
	"2006-01-02": {
		"2006/01/02",
		"20060102",
		"2006-01-02",
		"2006.01.02",
		"2006/1/2",
		"2006-1-2",
		"2006.1.2",
		"2 Jan 2006",
		"2Jan2006",
		"2 January 2006",
		"2January2006",
		"Jan 2, 2006",
		"January 2, 2006",
		// Explicitly prefer US Formatted MM/DD/YYYY
		"01/02/2006",
		"01-02-2006",
		"1/2/2006",
		"1-2-2006",
		"02/01/2006",
		"02-01-2006",
		"2/1/2006",
		"2-1-2006",
	},
	"2006-01": {
		"2006/01",
		"200601",
		"2006-01",
		"2006.01",
		"Jan 2006",
		"January 2006",
	},
	"2006": {"2006"},
}

func New(specifier string) (*DateOfBirth, error) {
	if specifier != "dateofbirth" {
		return nil, fmt.Errorf("invalid specifier: %s", specifier)
	}
	return &DateOfBirth{}, nil
}

func (dob *DateOfBirth) Apply(input string) (string, error) {
	result := ""
	for format, layouts := range supportedDateLayouts {
		for _, layout := range layouts {
			t, err := time.Parse(layout, input)
			if err == nil {
				if DateRangeCheck(t) {
					result = t.Format(format)
					break
				}
			}
		}
	}

	if len(result) == 0 {
		return "", fmt.Errorf("unrecognized date format")
	}

	return result, nil
}

func DateRangeCheck(t time.Time) bool {
	y, m, d := time.Now().Date()
	first := time.Date(y-120, m, d, 0, 0, 0, 0, time.UTC)
	last := time.Date(y, m, d+2, 0, 0, 0, 0, time.UTC)
	if t.After(first) && t.Before(last) {
		return true
	}
	return false
}
