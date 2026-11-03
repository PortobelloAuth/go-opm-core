package dateofbirth

import (
	"fmt"
	"time"
)

type DateOfBirth struct{}

var supportedDateLayouts = map[string][]string{
	"2006/01/02": []string{
		"2006/01/02",
		"20060102",
		"2006-01-02",
		"2006.01.02",
		"2 Jan 2006",
		"2Jan2006",
		"2 January 2006",
		"2January2006",
		"Jan 2, 2006",
		"January 2, 2006",
		"01/02/2006",
		"01-02-2006",
		"02/01/2006",
		"02-01-2006",
	},
	"2006/01": []string{
		"2006/01",
		"200601",
		"2006-01",
		"2006.01",
		"Jan 2006",
		"January 2006",
	},
	"2006": []string{"2006"},
}

func New(specifier string) (*DateOfBirth, error) {
	if specifier != "dateofbirth" {
		return nil, fmt.Errorf("invalid specifier: %s", specifier)
	}
	return &DateOfBirth{}, nil
}

func (l *DateOfBirth) Apply(input string) (string, error) {
	result := ""
	for format, layouts := range supportedDateLayouts {
		for _, layout := range layouts {
			t, err := time.Parse(layout, input)
			if err == nil {
				// TODO: make sure the date is within the expected range
				result = t.Format(format)
			}
		}
	}

	if len(result) == 0 {
		return "", fmt.Errorf("unrecognized date format")
	}

	return result, nil
}
