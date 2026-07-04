package normalization

import (
	"fmt"
	"strings"

	"github.com/PortobelloAuth/go-opm-core/lib/normalization/lowercase"
	"github.com/PortobelloAuth/go-opm-core/lib/normalization/replace"
	"github.com/PortobelloAuth/go-opm-core/lib/normalization/uppercase"
	"github.com/PortobelloAuth/go-opm-core/lib/normalization/whitespace"
)

type Normalization interface {
	Apply(input string) (string, error)
}

// TODO: consider whether we need short prefixes for all normalization specifiers
// TODO: figure out an interaction with rules that scores changes mode by a normalization (or
// treat changes that should be "scored" like replace as a distinct concept: transformations)

// New creates a new Normalization according to the specifier string
func New(specifier string) (Normalization, error) {
	switch specifier {
	case "lowercase":
		return lowercase.New(specifier)
	case "uppercase":
		return uppercase.New(specifier)
	case "whitespace":
		return whitespace.New(specifier)
	default:
		prefix, _, _ := strings.Cut(specifier, "/")
		switch prefix {
		case "s":
			return replace.New(specifier)
		}

		return nil, fmt.Errorf("%s is not a supported normalization", specifier)
	}
}
