package phonenumbere164

import (
	"fmt"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

var supportedRegions = phonenumbers.GetSupportedRegions()

type E164 struct {
	extensions    bool
	defaultRegion string
}

func New(specifier string) (*E164, error) {
	parts := strings.Split(specifier, "-")

	if len(parts) < 1 || parts[0] != "phonenumbere164" {
		return nil, fmt.Errorf("invalid specifier: %s", specifier)
	}

	normalizer := &E164{
		extensions:    false,
		defaultRegion: "US",
	}

	if len(parts) > 1 {
		if parts[1] == "extensions" {
			normalizer.extensions = true
			if len(parts) > 2 {
				normalizer.defaultRegion = parts[2]
			}
		} else {
			normalizer.defaultRegion = parts[1]
		}
	}

	if !supportedRegions[normalizer.defaultRegion] {
		return nil, fmt.Errorf("Unrecognized default region: %s", normalizer.defaultRegion)
	}

	return normalizer, nil
}

func (e164 *E164) Apply(input string) (string, error) {
	num, err := phonenumbers.Parse(input, e164.defaultRegion)
	if err != nil {
		return "", fmt.Errorf("failed to parse: %w", err)
	}

	// Validate against global numbering plans
	if !phonenumbers.IsValidNumber(num) {
		return "", fmt.Errorf("invalid phone number for parsed region")
	}

	formatted := phonenumbers.Format(num, phonenumbers.E164)
	ext := num.GetExtension()
	if e164.extensions && len(ext) > 0 {
		formatted = fmt.Sprintf("%s x%s", formatted, ext)
	}
	return formatted, nil
}
